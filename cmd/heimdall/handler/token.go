package handler

import (
	"errors"
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

var (
	BadRequestBodyError        = errors.New("failed to bind JSON body")
	TokenGenerationError       = errors.New("failed to generate token string")
	GetPayloadFromContextError = errors.New("failed get payload from context")
	PayloadTypeCastingError    = errors.New("failed to cast payload type")
	TokenRegexCreationError    = errors.New("failed to create token regex")
	TokenFormatError           = errors.New("token format is invalid")
	TokenParsingError          = errors.New("failed to parse token")
	TokenExpiredError          = errors.New("token is expired")
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
		_ = c.AbortWithError(http.StatusBadRequest, BadRequestBodyError)
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
		_ = c.AbortWithError(http.StatusInternalServerError, TokenGenerationError)
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
		_ = c.AbortWithError(http.StatusInternalServerError, GetPayloadFromContextError)
		return
	}
	payload, ok := payloadValue.(*config.Payload)
	if !ok {
		h.logger.Error("Failed parse payload type")
		_ = c.AbortWithError(http.StatusInternalServerError, PayloadTypeCastingError)
		return
	}
	c.JSON(http.StatusOK, payload)
}

func (h TokenHandler) VerifyAndSetHeader(c *gin.Context) {
	payloadValue, ok := c.Get("payload")
	if !ok {
		h.logger.Error("Failed get payload from context")
		_ = c.AbortWithError(http.StatusInternalServerError, GetPayloadFromContextError)
		return
	}
	payload, ok := payloadValue.(*config.Payload)
	if !ok {
		h.logger.Error("Failed parse payload type")
		_ = c.AbortWithError(http.StatusInternalServerError, PayloadTypeCastingError)
		return
	}

	for i := 0; i < reflect.TypeOf(payload.CustomPayload).NumField(); i++ {
		field := reflect.TypeOf(payload.CustomPayload).Field(i)
		headerName := field.Tag.Get("header")
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
		_ = c.AbortWithError(http.StatusInternalServerError, TokenRegexCreationError)
		return
	}
	if !reg.MatchString(bearerToken) {
		_ = c.AbortWithError(http.StatusUnauthorized, TokenFormatError)
		return
	}
	tokenString := reg.FindStringSubmatch(bearerToken)[1]
	payload, err := h.tokenManager.Parse(tokenString)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, TokenParsingError)
		return
	}

	if payload.ExpiredAt.Before(time.Now()) {
		_ = c.AbortWithError(http.StatusUnauthorized, TokenExpiredError)
		return
	}

	c.Set("payload", payload)
	c.Next()
}
