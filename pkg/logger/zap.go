package logger

import (
	"github.com/thetkpark/heimdall/pkg/config"
	"go.uber.org/zap"
)

func NewLogger(mode string) (logger *zap.Logger, err error) {
	if mode == config.ProductionMode {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
