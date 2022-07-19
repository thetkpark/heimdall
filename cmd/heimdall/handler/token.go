package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/token"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"time"
)

type TokenHandler struct {
	logger       *zap.SugaredLogger
	tokenManager token.Manager
	validTime    time.Duration
}

func NewTokenHandler(logger *zap.SugaredLogger, tokenMng token.Manager, validTime time.Duration) *TokenHandler {
	return &TokenHandler{
		logger:       logger,
		tokenManager: tokenMng,
		validTime:    validTime,
	}
}

func (h TokenHandler) GenerateToken(c *gin.Context) {
	var customPayload config.CustomPayload
	err := c.ShouldBindJSON(&customPayload)
	if err != nil {
		NewGinErrorResponse(c, http.StatusBadRequest, "Failed to bind JSON body")
		return
	}

	payload := config.Payload{
		CustomPayload: customPayload,
		MetadataPayload: config.MetadataPayload{
			IssuedAt:  time.Now().UTC(),
			ExpiredAt: time.Now().Add(h.validTime).UTC(),
		},
	}
	tokenString, err := h.tokenManager.Generate(payload)
	if err != nil {
		h.logger.Errorw("h.tokenManager.Generate error", "error", err, "payload", payload)
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed to generate token string")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"token": tokenString,
	})
}

func (h TokenHandler) VerifyToken(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	reg, err := regexp.Compile(`Bearer (.+\..+\..+)`)
	if err != nil {
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed to create regex against token")
		return
	}
	if !reg.MatchString(bearerToken) {
		NewGinErrorResponse(c, http.StatusBadRequest, "Token in Authorization header doesn't in the correct format")
		return
	}
	tokenString := reg.FindStringSubmatch(bearerToken)[1]
	payload, err := h.tokenManager.Parse(tokenString)
	if err != nil {
		h.logger.Errorw("h.tokenManager.Parse error", "error", err, "token", tokenString)
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed to verify token")
		return
	}

	if payload.ExpiredAt.Before(time.Now()) {
		NewGinErrorResponse(c, http.StatusUnauthorized, "Token is expired")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id": payload.UserID,
	})
}
