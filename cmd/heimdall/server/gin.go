package server

import (
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/thetkpark/heimdall/cmd/heimdall/handler"
	"github.com/thetkpark/heimdall/pkg/config"
	"net/http"
	"time"
)

func NewGINServer(cfg *config.Config, tokenHandler *handler.TokenHandler) *http.Server {
	gin.SetMode(cfg.GinMode)
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	router.Use(handler.HTTPErrorHandler)
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})
	router.GET("/verify", tokenHandler.AuthenticateToken, tokenHandler.VerifyToken)
	router.GET("/auth", tokenHandler.AuthenticateToken, tokenHandler.VerifyAndSetHeader)
	router.POST("/generate", tokenHandler.GenerateToken)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.GinPort),
		Handler: router,
	}

	return httpServer
}
