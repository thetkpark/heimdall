package config

import "time"

type MetadataPayload struct {
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CustomPayload struct {
	UserID uint32 `json:"user_id" binding:"required" header:"X-USER-ID"`
}

type Payload struct {
	CustomPayload
	MetadataPayload
}
