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
	SentryDSN            string        `env:"SENTRY_DSN"`
	Mode                 string        `env:"MODE" envDefault:"development"`
	GinMode              string        `env:"GIN_MODE" envDefault:"debug"`
	GinPort              int           `env:"GIN_PORT" envDefault:"8080"`
	GRPCPort             int           `env:"GRPC_PORT" envDefault:"5050"`
}

func ParseConfig() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
