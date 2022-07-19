package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	grpc2 "github.com/thetkpark/heimdall/cmd/heimdall/grpc"
	"github.com/thetkpark/heimdall/cmd/heimdall/handler"
	pb "github.com/thetkpark/heimdall/cmd/heimdall/proto"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/logger"
	"github.com/thetkpark/heimdall/pkg/signature"
	"github.com/thetkpark/heimdall/pkg/token"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
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
		Dsn:              "https://00b3b70ef73a48dfa377bad2fe9af84a@o1318116.ingest.sentry.io/6572271",
		TracesSampleRate: 1.0,
	}); err != nil {
		sugaredLogger.Fatalw("Failed to init Sentry", "error", err)
	}
	defer sentry.Flush(2 * time.Second)

	ginLogger := sugaredLogger.Named("gin")
	gin.SetMode(cfg.GinMode)
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})
	signatureManager := signature.NewJWS(cfg.JWSSecretKey)
	tokenManager := token.NewTokenManager(signatureManager, nil)
	tokenHandler := handler.NewTokenHandler(ginLogger, tokenManager, cfg.TokenValidTime)
	router.GET("/verify", tokenHandler.AuthenticateToken, tokenHandler.VerifyToken)
	router.GET("/auth", tokenHandler.AuthenticateToken, tokenHandler.VerifyAndSetHeader)
	router.POST("/generate", tokenHandler.GenerateToken)
	go endless.ListenAndServe(":8080", router)

	// gRPC
	grpcLogger := sugaredLogger.Named("grpc")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5050))
	if err != nil {
		grpcLogger.Fatalw("Failed to listen", "error", err, "port", 5050)
	}

	grpcServer := grpc.NewServer()
	grpcTokenServer := grpc2.NewTokenServer(grpcLogger, tokenManager, cfg.TokenValidTime)
	pb.RegisterTokenServer(grpcServer, grpcTokenServer)
	grpcLogger.Infof("Starting gRPC server on port")
	if err := grpcServer.Serve(lis); err != nil {
		grpcLogger.Fatalw("Failed to start gRPC server", "error", err)
	}
}
