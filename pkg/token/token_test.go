package token_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/token"
	"github.com/thetkpark/heimdall/test/mock_encryption"
	"github.com/thetkpark/heimdall/test/mock_signature"
	"time"
)

var _ = Describe("Token Manager", func() {
	var (
		mockCtrl       *gomock.Controller
		mockEncryption *mock_encryption.MockManager
		mockSignature  *mock_signature.MockManager
		tokenManager   token.Manager
		payload        config.Payload
		rawPayload     []byte
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEncryption = mock_encryption.NewMockManager(mockCtrl)
		mockSignature = mock_signature.NewMockManager(mockCtrl)
		payload = config.Payload{
			CustomPayload: config.CustomPayload{UserID: 99},
			MetadataPayload: config.MetadataPayload{
				IssuedAt:  time.Unix(1658201439, 0),
				ExpiredAt: time.Unix(1658201439, 0).Add(time.Second * 10),
			},
		}
		var err error
		rawPayload, err = json.Marshal(payload)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Unencrypted payload token", func() {
		BeforeEach(func() {
			tokenManager = token.NewTokenManager(mockSignature, nil)
		})

		It("can generate the token", func() {
			mockSignedToken := []byte("signedToken")
			mockSignature.EXPECT().Sign(rawPayload).Return(mockSignedToken, nil).Times(1)
			token, err := tokenManager.Generate(payload)
			Expect(err).To(BeNil())
			Expect([]byte(token)).To(Equal(mockSignedToken))
		})

		It("can parse the token", func() {
			mockSignedToken := "signedToken"
			mockSignature.EXPECT().Verify([]byte(mockSignedToken)).Return(rawPayload, nil).Times(1)
			retrievedPayload, err := tokenManager.Parse(mockSignedToken)
			Expect(err).To(BeNil())
			Expect(*retrievedPayload).To(Equal(payload))
		})
	})

	Context("Encrypted payload token", func() {
		BeforeEach(func() {
			tokenManager = token.NewTokenManager(mockSignature, mockEncryption)
		})

		It("can generate the token", func() {
			mockSignedToken := []byte("signedToken")
			mockEncryptedRawPayload := []byte("encryptedRawPayload")
			mockEncryption.EXPECT().Encrypt(rawPayload).Return(mockEncryptedRawPayload, nil).Times(1)
			mockSignature.EXPECT().Sign(mockEncryptedRawPayload).Return(mockSignedToken, nil).Times(1)

			token, err := tokenManager.Generate(payload)
			Expect(err).To(BeNil())
			Expect([]byte(token)).To(Equal(mockSignedToken))
		})

		It("can parse the token", func() {
			mockSignedToken := "signedToken"
			mockEncryptedRawPayload := []byte("encryptedRawPayload")
			mockSignature.EXPECT().Verify([]byte(mockSignedToken)).Return(mockEncryptedRawPayload, nil).Times(1)
			mockEncryption.EXPECT().Decrypt(mockEncryptedRawPayload).Return(rawPayload, nil).Times(1)
			retrievedPayload, err := tokenManager.Parse(mockSignedToken)
			Expect(err).To(BeNil())
			Expect(*retrievedPayload).To(Equal(payload))
		})
	})
})
