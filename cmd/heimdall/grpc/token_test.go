package grpc_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thetkpark/heimdall/cmd/heimdall/grpc"
	pb "github.com/thetkpark/heimdall/cmd/heimdall/proto"
	"github.com/thetkpark/heimdall/test/mock_token"
	"go.uber.org/zap"
	"time"
)

var _ = Describe("TokenServer_gRPC", func() {

	var (
		mockCtrl         *gomock.Controller
		mockTokenManager *mock_token.MockManager
		tokenServer      *grpc.TokenServer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTokenManager = mock_token.NewMockManager(mockCtrl)
		tokenServer = grpc.NewTokenServer(zap.NewNop().Sugar(), mockTokenManager, time.Hour)
	})

	Context("GenerateToken", func() {
		var (
			handler  func(_ context.Context, tokenReq *pb.GenerateTokenRequest) (*pb.TokenResponse, error)
			req      *pb.GenerateTokenRequest
			res      *pb.TokenResponse
			resError error
		)

		BeforeEach(func() {
			handler = tokenServer.GenerateToken
			req = &pb.GenerateTokenRequest{UserID: 99999}
		})

		JustBeforeEach(func() {
			res, resError = handler(context.Background(), req)
		})

		When("Request is valid", func() {
			BeforeEach(func() {
				mockTokenManager.EXPECT().Generate(gomock.Any()).Return("token", nil).Times(1)
			})

			It("should get the token successfully", func() {
				Expect(resError).To(BeNil())
				Expect(res.Token).To(Equal("token"))
			})
		})

		When("Failed to generate token", func() {
			BeforeEach(func() {
				mockTokenManager.EXPECT().Generate(gomock.Any()).Return("", errors.New("failed to generate")).Times(1)
			})

			It("should return Internal error", func() {
				Expect(resError).ToNot(BeNil())
			})
		})

	})
})
