package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to intialized Zap Logger: %v", err.Error())
	}
	defer zapLogger.Sync()
	//logger := zapLogger.Sugar()

	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	_ = endless.ListenAndServe(":8080", router)
}
