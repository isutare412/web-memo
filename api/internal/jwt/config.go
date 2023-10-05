package jwt

import (
	"time"
)

type Config struct {
	ActiveKeyPair string            `mapstructure:"active-key-pair" validate:"required"`
	KeyPairs      []RSAKeyBytesPair `mapstructure:"key-pairs" validate:"required"`
	Expiration    time.Duration     `mapstructure:"expiration" validate:"required"`
}

type RSAKeyBytesPair struct {
	Name    string `mapstrucutre:"name" validate:"required"`
	Private string `mapstrucutre:"private" validate:"required"`
	Public  string `mapstrucutre:"public" validate:"required"`
}
