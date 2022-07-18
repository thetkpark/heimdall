package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type Encryption interface {
	Encrypt(plainText string) ([]byte, error)
	Decrypt(cipherText []byte) (string, error)
}

func NewAESEncryption(key []byte) (*AES, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	return &AES{
		gcm: gcm,
	}, nil
}

type AES struct {
	gcm cipher.AEAD
}

func (a AES) Encrypt(plaintext string) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherText := a.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return cipherText, nil
}

func (a AES) Decrypt(cipherText []byte) (string, error) {
	nonceSize := a.gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext length is shorter than nonce size")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := a.gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
