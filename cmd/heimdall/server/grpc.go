package server

import (
	grpc2 "github.com/thetkpark/heimdall/cmd/heimdall/grpc"
	pb "github.com/thetkpark/heimdall/cmd/heimdall/proto"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/token"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewGRPCServer(logger *zap.SugaredLogger, cfg *config.Config, tokenMng token.Manager) *grpc.Server {
	grpcServer := grpc.NewServer()
	grpcTokenServer := grpc2.NewTokenServer(logger, tokenMng, cfg.TokenValidTime)
	pb.RegisterTokenServer(grpcServer, grpcTokenServer)
	return grpcServer
}
