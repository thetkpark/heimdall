package config

import "time"

type MetadataPayload struct {
	IssuedAt  time.Time
	ExpiredAt time.Time
}

type CustomPayload struct {
	UserID uint32 `header:"X-USER_ID"`
}

type Payload struct {
	CustomPayload
	MetadataPayload
}
