package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseConfig = `
process-timeout: 1m
aws:
  access-key: test_key
  secret: test_secret
  s3:
    backup-bucket: test_bucket
backup:
  prefix: some-prefix/
  retention: 6h
postgres:
  host: localhost
  port: 8080
  user: test_user
  password: test_pass
  database: test_database
`

func TestLoadValidated(t *testing.T) {
	type args struct {
		config      string
		localConfig string
		envs        map[string]string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "only config",
			args: args{
				config: baseConfig,
			},
			want: &Config{
				ProcessTimeout: time.Minute,
				AWS: AWSConfig{
					AccessKey: "test_key",
					Secret:    "test_secret",
					S3: S3Config{
						BackupBucket: "test_bucket",
					},
				},
				Backup: BackupConfig{
					Prefix:    "some-prefix/",
					Retention: 6 * time.Hour,
				},
				Postgres: PostgresConfig{
					Host:     "localhost",
					Port:     8080,
					User:     "test_user",
					Password: "test_pass",
					Database: "test_database",
				},
			},
		},
		{
			name: "config overwitten by local config",
			args: args{
				config: baseConfig,
				localConfig: `
aws:
  s3:
    backup-bucket: other_bucket
`,
			},
			want: &Config{
				ProcessTimeout: time.Minute,
				AWS: AWSConfig{
					AccessKey: "test_key",
					Secret:    "test_secret",
					S3: S3Config{
						BackupBucket: "other_bucket",
					},
				},
				Backup: BackupConfig{
					Prefix:    "some-prefix/",
					Retention: 6 * time.Hour,
				},
				Postgres: PostgresConfig{
					Host:     "localhost",
					Port:     8080,
					User:     "test_user",
					Password: "test_pass",
					Database: "test_database",
				},
			},
		},
		{
			name: "config overwitten by envs",
			args: args{
				config: baseConfig,
				envs: map[string]string{
					"APP_AWS_S3_BACKUP-BUCKET": "another_bucket",
				},
			},
			want: &Config{
				ProcessTimeout: time.Minute,
				AWS: AWSConfig{
					AccessKey: "test_key",
					Secret:    "test_secret",
					S3: S3Config{
						BackupBucket: "another_bucket",
					},
				},
				Backup: BackupConfig{
					Prefix:    "some-prefix/",
					Retention: 6 * time.Hour,
				},
				Postgres: PostgresConfig{
					Host:     "localhost",
					Port:     8080,
					User:     "test_user",
					Password: "test_pass",
					Database: "test_database",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()

			err := os.WriteFile(filepath.Join(testDir, configFileNames[0]), []byte(tt.args.config), 0644)
			require.NoError(t, err)

			if tt.args.localConfig != "" {
				err := os.WriteFile(filepath.Join(testDir, configFileNames[1]), []byte(tt.args.localConfig), 0644)
				require.NoError(t, err)
			}

			for k, v := range tt.args.envs {
				t.Setenv(k, v)
			}

			got, err := LoadValidated(testDir)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
