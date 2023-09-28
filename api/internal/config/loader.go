package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/isutare412/tasks/api/internal/validate"
)

func LoadValidated(path string) (*Config, error) {
	if err := readFile(path); err != nil {
		return nil, err
	}
	readEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, nil
	}

	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validating loaded config: %w", err)
	}
	return &cfg, nil
}

func readFile(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("reading in config: %w", err)
	}
	return nil
}

func readEnv() {
	// APP_FOO_BAR=baz -> cfg.Foo.Bar = "baz"
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
