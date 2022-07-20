package token

import (
	"encoding/json"
	"github.com/thetkpark/heimdall/pkg/config"
	"github.com/thetkpark/heimdall/pkg/encryption"
	"github.com/thetkpark/heimdall/pkg/signature"
)

type Manager interface {
	Generate(payload config.Payload) (string, error)
	Parse(token string) (*config.Payload, error)
}

func NewTokenManager(sig signature.Manager, enc encryption.Manager) *manager {
	mng := &manager{encryptionManager: enc, signatureManager: sig}
	return mng
}

type manager struct {
	signatureManager  signature.Manager
	encryptionManager encryption.Manager
}

func (m *manager) SetEncryptionManager(enc encryption.Manager) {
	m.encryptionManager = enc
}

func (m manager) Generate(payload config.Payload) (string, error) {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	if m.encryptionManager != nil {
		rawPayload, err = m.encryptionManager.Encrypt(rawPayload)
		if err != nil {
			return "", err
		}
	}

	token, err := m.signatureManager.Sign(rawPayload)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func (m manager) Parse(token string) (*config.Payload, error) {
	rawPayload, err := m.signatureManager.Verify([]byte(token))
	if err != nil {
		return nil, err
	}

	if m.encryptionManager != nil {
		rawPayload, err = m.encryptionManager.Decrypt(rawPayload)
		if err != nil {
			return nil, err
		}
	}

	var payload config.Payload
	err = json.Unmarshal(rawPayload, &payload)
	return &payload, err
}
