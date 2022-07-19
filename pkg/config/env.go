package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"time"
)

const ProductionMode = "production"

type Config struct {
	JWSSecretKey         string        `env:"JWS_SECRET_KEY,required"`
	PayloadEncryptionKey string        `env:"PAYLOAD_ENCRYPTION_KEY"`
	TokenValidTime       time.Duration `env:"TOKEN_VALID_TIME"`
	Mode                 string        `env:"MODE" envDefault:"development"`
	GinMode              string        `env:"GIN_MODE" envDefault:"debug"`
}

func ParseConfig() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
