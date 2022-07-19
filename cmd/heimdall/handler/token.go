package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/token"
	"go.uber.org/zap"
	"net/http"
	"reflect"
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
	payloadValue, ok := c.Get("payload")
	if !ok {
		h.logger.Error("Failed get payload from context")
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed get payload from context")
		return
	}
	payload, ok := payloadValue.(*config.Payload)
	if !ok {
		h.logger.Error("Failed parse payload type")
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed parse payload type")
		return
	}
	c.JSON(http.StatusOK, payload)
}

func (h TokenHandler) VerifyAndSetHeader(c *gin.Context) {
	payloadValue, ok := c.Get("payload")
	if !ok {
		h.logger.Error("Failed get payload from context")
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed get payload from context")
		return
	}
	payload, ok := payloadValue.(*config.Payload)
	if !ok {
		h.logger.Error("Failed parse payload type")
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed parse payload type")
		return
	}

	for i := 0; i < reflect.TypeOf(payload.CustomPayload).NumField(); i++ {
		field := reflect.TypeOf(payload.CustomPayload).Field(i)
		headerName := field.Tag.Get("header")
		h.logger.Info("headerName", headerName)
		if len(headerName) > 0 {
			val := fmt.Sprintf("%v", reflect.ValueOf(payload.CustomPayload).Field(i))
			c.Header(headerName, val)
		}
	}
	c.Status(http.StatusOK)
}

func (h TokenHandler) AuthenticateToken(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	reg, err := regexp.Compile(`Bearer (.+\..+\..+)`)
	if err != nil {
		NewGinErrorResponse(c, http.StatusInternalServerError, "Failed to create regex against token")
		return
	}
	if !reg.MatchString(bearerToken) {
		NewGinErrorResponse(c, http.StatusUnauthorized, "Token in Authorization header doesn't in the correct format")
		return
	}
	tokenString := reg.FindStringSubmatch(bearerToken)[1]
	payload, err := h.tokenManager.Parse(tokenString)
	if err != nil {
		NewGinErrorResponse(c, http.StatusUnauthorized, "Token verification failed")
		return
	}

	if payload.ExpiredAt.Before(time.Now()) {
		NewGinErrorResponse(c, http.StatusUnauthorized, "Token is expired")
		return
	}

	c.Set("payload", payload)
	c.Next()
}
