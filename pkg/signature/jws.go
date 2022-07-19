package signature

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
	goJWS "github.com/lestrrat-go/jwx/v2/jws"
)

type Manager interface {
	Sign(payload []byte) ([]byte, error)
	Verify(token []byte) ([]byte, error)
}

type jws struct {
	encryptionKey []byte
}

func NewJWS(key string) *jws {
	return &jws{encryptionKey: []byte(key)}
}

func (j jws) Sign(payload []byte) ([]byte, error) {
	return goJWS.Sign(payload, goJWS.WithKey(jwa.HS256, j.encryptionKey))
}

func (j jws) Verify(token []byte) ([]byte, error) {
	return goJWS.Verify(token, goJWS.WithKey(jwa.HS256, j.encryptionKey))
}
