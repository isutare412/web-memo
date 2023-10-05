package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/isutare412/web-memo/api/internal/validate"
)

var overrideFiles = []string{
	"config.local.yaml",
}

func LoadValidated(path string) (*Config, error) {
	if err := readFiles(path); err != nil {
		return nil, err
	}
	readEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validating loaded config: %w", err)
	}
	return &cfg, nil
}

func readFiles(dirPath string) error {
	viper.SetConfigFile(filepath.Join(dirPath, "config.yaml"))
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("reading in config: %w", err)
	}

	for _, fileName := range overrideFiles {
		if _, err := mergeFileIfExist(filepath.Join(dirPath, fileName)); err != nil {
			return err
		}
	}
	return nil
}

func mergeFileIfExist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("checking if config exists: %w", err)
	}

	viper.SetConfigFile(path)
	if err := viper.MergeInConfig(); err != nil {
		return false, fmt.Errorf("merging in config: %w", err)
	}

	return true, nil
}

func readEnv() {
	// APP_FOO_BAR=baz -> cfg.Foo.Bar = "baz"
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
