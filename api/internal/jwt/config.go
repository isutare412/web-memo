package jwt

import (
	"time"
)

type Config struct {
	ActiveKeyPair string            `koanf:"active-key-pair" validate:"required"`
	KeyPairs      []RSAKeyBytesPair `koanf:"key-pairs" validate:"required"`
	Expiration    time.Duration     `koanf:"expiration" validate:"required"`
}

type RSAKeyBytesPair struct {
	Name    string `koanf:"name" validate:"required"`
	Private string `koanf:"private" validate:"required"`
	Public  string `koanf:"public" validate:"required"`
}
