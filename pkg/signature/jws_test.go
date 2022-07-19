package signature_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thetkpark/heimdall/pkg/signature"
)

var _ = Describe("Json Web Signature", func() {
	var jws signature.Manager
	key := "E2sK$Cps7v1sB2RW010HlSWdpS&CSOy4"
	plaintext := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua")

	BeforeEach(func() {
		jws = signature.NewJWS(key)
	})

	It("can sign and verify signature", func() {
		token, err := jws.Sign(plaintext)
		Expect(err).To(BeNil())
		Expect(token).ToNot(BeEmpty())

		verified, err := jws.Verify(token)
		Expect(err).To(BeNil())
		Expect(verified).To(Equal(plaintext))
	})

	It("failed to verify when token is edited", func() {
		token := "eyJhbGciOiJIUzI1NiJ9.TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnRldHVyIGFkaXBpc2NpbmcgZWxpdCwgc2VkIGRvIGVpdXNtb2QgdGVtcG9yIGluY2lkaWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdWE.Sp4CSEJgAe_AS2Ao2cyO3K17ufOGieZTpjzwrKwHl7U\n\n"
		verified, err := jws.Verify([]byte(token))
		Expect(err).ToNot(BeNil())
		Expect(verified).To(BeNil())
	})
})
