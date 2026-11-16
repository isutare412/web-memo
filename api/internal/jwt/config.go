package jwt

import (
	"time"
)

type Config struct {
	ActiveKeyPair string `koanf:"active-key-pair" validate:"required"`
	// KeyPairs maps a key pair name, which tokens carry in the kid header, to its keys.
	// Keyed by name so each key can also come from env, e.g. APP_JWT_KEY-PAIRS_V1_PRIVATE;
	// env keys are lowercased and "_" splits levels, so names must be lowercase without "_".
	KeyPairs   map[string]RSAKeyBytesPair `koanf:"key-pairs" validate:"required,dive"`
	Expiration time.Duration              `koanf:"expiration" validate:"required"`
}

type RSAKeyBytesPair struct {
	Private string `koanf:"private" validate:"required"`
	Public  string `koanf:"public" validate:"required"`
}
