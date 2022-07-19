package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/logger"
	"log"
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
	//logger := zapLogger.Sugar()

	gin.SetMode(cfg.GinMode)
	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	_ = endless.ListenAndServe(":8080", router)
}
