package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/thetkpark/heimdall/cmd/heimdall/handler"
	"github.com/thetkpark/heimdall/cmd/heimdall/server"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/encryption"
	"github.com/thetkpark/heimdall/pkg/logger"
	"github.com/thetkpark/heimdall/pkg/signature"
	"github.com/thetkpark/heimdall/pkg/token"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Failed to parse ENV: %v", err)
	}

	zapLogger, err := logger.NewLogger(cfg.Mode)
	if err != nil {
		log.Fatalf("Failed to initialized Zap: %v", err)
	}
	defer zapLogger.Sync()
	sugaredLogger := zapLogger.Sugar()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		TracesSampleRate: 1.0,
	}); err != nil {
		sugaredLogger.Fatalw("Failed to init Sentry", "error", err)
	}
	defer sentry.Flush(3 * time.Second)

	signatureManager := signature.NewJWS(cfg.JWSSecretKey)
	tokenManager := token.NewTokenManager(signatureManager, nil)
	if len(cfg.PayloadEncryptionKey) > 0 {
		encryptionManager, err := encryption.NewAESEncryption([]byte(cfg.PayloadEncryptionKey))
		if err != nil {
			sugaredLogger.Fatalw("Failed to init AESEncryption", "error", err)
		}
		tokenManager.SetEncryptionManager(encryptionManager)
	}
	tokenHandler := handler.NewTokenHandler(sugaredLogger, tokenManager, cfg.TokenValidTime)

	ginLogger := sugaredLogger.Named("GIN")
	ginServer := server.NewGINServer(cfg, tokenHandler)
	go func() {
		ginLogger.Infof("Starting GIN server on %d", cfg.GinPort)
		if err := ginServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			ginLogger.Infof("GIN HTTP listen error: %v\n", err)
		}
	}()

	// gRPC
	grpcLogger := sugaredLogger.Named("gRPC")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		grpcLogger.Fatalw("Failed to listen", "error", err, "port", 5050)
	}
	grpcServer := server.NewGRPCServer(grpcLogger, cfg, tokenManager)
	go func() {
		grpcLogger.Infof("Starting gRPC server on %d", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			grpcLogger.Fatalw("Failed to start gRPC server", "error", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugaredLogger.Info("SIG received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := ginServer.Shutdown(ctx); err != nil {
		ginLogger.Fatal("GIN server forced to shutdown:", err)
	}

	sugaredLogger.Info("Server exiting")
}
