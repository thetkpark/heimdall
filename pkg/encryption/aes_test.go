package encryption_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thetkpark/heimdall/pkg/encryption"
)

var _ = Describe("AES Encryption", Label("encryption"), func() {
	var aes *encryption.AES
	plaintext := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	BeforeEach(func() {
		var err error
		aes, err = encryption.NewAESEncryption([]byte("E2sK$Cps7v1sB2RW010HlSWdpS&CSOy4"))
		Expect(err).To(BeNil())
	})

	It("can encrypt and decrypt", func() {
		base64CipherText, err := aes.Encrypt(plaintext)
		Expect(err).To(BeNil())
		Expect(base64CipherText).ToNot(BeEmpty())

		decryptedPlainText, err := aes.Decrypt(base64CipherText)
		Expect(err).To(BeNil())
		Expect(decryptedPlainText).To(Equal(plaintext))
	})
})
