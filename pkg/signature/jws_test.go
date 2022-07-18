package signature_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thetkpark/heimdall/pkg/signature"
)

var _ = Describe("Json Web Signature", func() {
	var jws signature.Signature
	key := "E2sK$Cps7v1sB2RW010HlSWdpS&CSOy4"
	plaintext := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua")

	BeforeEach(func() {
		jws = signature.NewJWS(key)
	})

	It("can sign and verify signature", func() {
		token, err := jws.Sign(plaintext)
		Expect(err).To(BeNil())
		Expect(token).ToNot(BeEmpty())
		fmt.Println(string(token))

		verified, err := jws.Verify(token)
		Expect(err).To(BeNil())
		Expect(verified).To(Equal(plaintext))
		fmt.Println(string(verified))
	})
})
