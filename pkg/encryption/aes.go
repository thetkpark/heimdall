package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Encryption interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
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

func (a AES) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := a.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a AES) Decrypt(base64CipherText string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return "", err
	}

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
