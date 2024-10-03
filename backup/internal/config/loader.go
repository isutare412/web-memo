package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"

	"github.com/isutare412/web-memo/backup/internal/validate"
)

var configFileNames = [...]string{
	"config.yaml",
	"config.local.yaml",
}

func LoadValidated(dir string) (*Config, error) {
	k := koanf.New(".")

	var configFiles []string
	for _, name := range configFileNames {
		configFiles = append(configFiles, filepath.Join(dir, name))
	}

	for i, f := range configFiles {
		if i != 0 { // file is optional except first one
			if _, err := os.Stat(f); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return nil, fmt.Errorf("checking config file existence: %w", err)
			}
		}

		if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("loading from file(%s): %w", f, err)
		}
	}

	// APP_FOO_BAR=baz -> foo.bar=baz
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "APP_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("loading from env: %w", err)
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, fmt.Errorf("unmarshaling into config struct: %w", err)
	}
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validating loaded config: %w", err)
	}

	return &cfg, nil
}
