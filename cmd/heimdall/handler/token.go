package handler

import (
	"errors"
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
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

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewTokenHandler(logger *zap.SugaredLogger, tokenMng token.Manager, validTime time.Duration) *TokenHandler {
	return &TokenHandler{
		logger:       logger,
		tokenManager: tokenMng,
		validTime:    validTime,
	}
}

// GenerateToken godoc
// @Summary      Generate token with the payload
// @Tags         token
// @Accept       json
// @Produce      json
// @Param payload body config.CustomPayload true "Payload"
// @Success      201  {object}  TokenResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /generate [POST]
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
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(err)
		}
		return
	}
	c.JSON(http.StatusCreated, TokenResponse{Token: tokenString})
}

// ParsePayload godoc
// @Summary      Verify token and parse payload
// @Tags         token
// @Security	 JWSToken
// @Produce      json
// @Success      200  {object}  config.Payload
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/body [GET]
func (h TokenHandler) ParsePayload(c *gin.Context) {
	payloadValue, ok := c.Get("payload")
	if !ok {
		h.logger.Error("Failed get payload from context")
		_ = c.AbortWithError(http.StatusInternalServerError, GetPayloadFromContextError)
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(GetPayloadFromContextError)
		}
		return
	}
	payload, ok := payloadValue.(*config.Payload)
	if !ok {
		h.logger.Error("Failed parse payload type")
		_ = c.AbortWithError(http.StatusInternalServerError, PayloadTypeCastingError)
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(PayloadTypeCastingError)
		}
		return
	}
	c.JSON(http.StatusOK, payload)
}

// ParsePayloadAndSetHeader godoc
// @Summary      Verify token and set custom payload to header
// @Tags         token
// @Security	 JWSToken
// @Success      200
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/header [GET]
func (h TokenHandler) ParsePayloadAndSetHeader(c *gin.Context) {
	payloadValue, ok := c.Get("payload")
	if !ok {
		h.logger.Error("Failed get payload from context")
		_ = c.AbortWithError(http.StatusInternalServerError, GetPayloadFromContextError)
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(GetPayloadFromContextError)
		}
		return
	}
	payload, ok := payloadValue.(*config.Payload)
	if !ok {
		h.logger.Error("Failed parse payload type")
		_ = c.AbortWithError(http.StatusInternalServerError, PayloadTypeCastingError)
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(PayloadTypeCastingError)
		}
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
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(err)
		}
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

	if h.isTokenExpired(payload.ExpiredAt) {
		_ = c.AbortWithError(http.StatusUnauthorized, TokenExpiredError)
		return
	}

	c.Set("payload", payload)
	c.Next()
}

func (h TokenHandler) isTokenExpired(expiredAt time.Time) bool {
	if h.validTime.Microseconds() == 0 || expiredAt.After(time.Now()) {
		return false
	}
	return true
}
