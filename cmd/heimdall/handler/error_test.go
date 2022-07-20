package handler_test

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thetkpark/heimdall/cmd/heimdall/handler"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ErrorHandler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		handlerFunc gin.HandlerFunc
		errorStatus int
		testError   error
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		rec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(rec)
		handlerFunc = handler.HTTPErrorHandler

		testError = errors.New("test error")
		errorStatus = http.StatusTeapot
		_ = c.AbortWithError(errorStatus, testError)
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should return error in JSON format", func() {
		Expect(rec.Body.String()).To(Equal(fmt.Sprintf(`{"error":"%s"}`, testError.Error())))
	})

	It("should not manipulate the status", func() {
		Expect(rec.Code).To(Equal(errorStatus))
	})
})
