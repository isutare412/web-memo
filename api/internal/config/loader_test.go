package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isutare412/tasks/api/internal/log"
)

func TestLoadValidated(t *testing.T) {
	tests := []struct {
		name      string
		rawConfig string
		envs      map[string]string
		want      *Config
		wantErr   bool
	}{
		{
			name: "normal_case",
			rawConfig: `log:
  format: text
  level: debug
  caller: true
postgres:
  host: 127.0.0.1
  port: 1234
  user: tester
  password: password
  database: fake`,
			want: &Config{
				Log: LogConfig{
					Format: log.FormatText,
					Level:  log.LevelDebug,
					Caller: true,
				},
				Postgres: PostgresConfig{
					Host:     "127.0.0.1",
					Port:     1234,
					User:     "tester",
					Password: "password",
					Database: "fake",
				},
			},
		},
		{
			name: "overriden_by_env",
			rawConfig: `log:
  format: text
  level: debug
  caller: true
postgres:
  host: 127.0.0.1
  port: 1234
  user: tester
  password: password
  database: fake`,
			envs: map[string]string{
				"APP_POSTGRES_HOST": "1.2.3.4",
			},
			want: &Config{
				Log: LogConfig{
					Format: log.FormatText,
					Level:  log.LevelDebug,
					Caller: true,
				},
				Postgres: PostgresConfig{
					Host:     "1.2.3.4",
					Port:     1234,
					User:     "tester",
					Password: "password",
					Database: "fake",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(file, []byte(tt.rawConfig), 0644); err != nil {
				require.NoError(t, err)
			}

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			got, err := LoadValidated(file)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
