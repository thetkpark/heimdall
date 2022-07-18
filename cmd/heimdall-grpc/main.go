package main

import (
	"fmt"
	pb "github.com/thetkpark/heimdall/cmd/heimdall-grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedTokenServer
}

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to intialized Zap Logger: %v", err.Error())
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		logger.Fatalw("Failed to listen", "error", err, "port", 8080)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTokenServer(grpcServer, &server{})
	logger.Info("Starting gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalw("Failed to start gRPC server", "error", err)
	}
}
