package config_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/config"
	"github.com/isutare412/web-memo/api/internal/http"
	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

var _ = Describe("Loader", func() {
	Context("LoadValidated", func() {
		var (
			givenConfigDir  string
			givenConfigText = `wire:
  initialize-timeout: 1h30m
  shutdown-timeout: 90s
log:
  format: text
  level: debug
  caller: true
http:
  port: 8412
postgres:
  host: 127.0.0.1
  port: 1234
  user: tester
  password: password
  database: fake
redis:
  addr: localhost:4120
  password: memcached
google:
  endpoints:
    token: https://oauth2.googleapis.com/token
    oauth: https://accounts.google.com/o/oauth2/v2/auth
  oauth:
    client-id: google-oauth-client-id
    client-secret: google-oauth-client-secret
service:
  auth:
    oauth-state-timeout: 42m
    google-callback-path: /api/v1/google/sign-in/finish`
		)

		BeforeEach(func() {
			givenConfigDir = GinkgoT().TempDir()

			file := filepath.Join(givenConfigDir, "config.yaml")
			err := os.WriteFile(file, []byte(givenConfigText), 0644)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("loads valid config", func() {
			var (
				wantConfig = config.Config{
					Wire: config.WireConfig{
						InitializeTimeout: 90 * time.Minute,
						ShutdownTimeout:   90 * time.Second,
					},
					Log: log.Config{
						Format: log.FormatText,
						Level:  log.LevelDebug,
						Caller: true,
					},
					HTTP: http.Config{
						Port: 8412,
					},
					Postgres: postgres.Config{
						Host:     "127.0.0.1",
						Port:     1234,
						User:     "tester",
						Password: "password",
						Database: "fake",
					},
					Redis: redis.Config{
						Addr:     "localhost:4120",
						Password: "memcached",
					},
					Google: config.GoogleConfig{
						Endpoints: config.GoogleEndpointsConfig{
							Token: "https://oauth2.googleapis.com/token",
							OAuth: "https://accounts.google.com/o/oauth2/v2/auth",
						},
						OAuth: config.GoogleOAuthConfig{
							ClientID:     "google-oauth-client-id",
							ClientSecret: "google-oauth-client-secret",
						},
					},
					Service: config.ServiceConfig{
						Auth: config.AuthServiceConfig{
							OAuthStateTimeout:  42 * time.Minute,
							GoogleCallbackPath: "/api/v1/google/sign-in/finish",
						},
					},
				}
			)

			cfg, err := config.LoadValidated(givenConfigDir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*cfg).Should(Equal(wantConfig))
		})

		It("loads config overwritten by env", func() {
			var (
				givenEnvs = map[string]string{
					"APP_POSTGRES_HOST":         "1.2.3.4",
					"APP_WIRE_SHUTDOWN-TIMEOUT": "1s",
				}
			)

			var (
				wantConfig = config.Config{
					Wire: config.WireConfig{
						InitializeTimeout: 90 * time.Minute,
						ShutdownTimeout:   time.Second,
					},
					Log: log.Config{
						Format: log.FormatText,
						Level:  log.LevelDebug,
						Caller: true,
					},
					HTTP: http.Config{
						Port: 8412,
					},
					Postgres: postgres.Config{
						Host:     "1.2.3.4",
						Port:     1234,
						User:     "tester",
						Password: "password",
						Database: "fake",
					},
					Redis: redis.Config{
						Addr:     "localhost:4120",
						Password: "memcached",
					},
					Google: config.GoogleConfig{
						Endpoints: config.GoogleEndpointsConfig{
							Token: "https://oauth2.googleapis.com/token",
							OAuth: "https://accounts.google.com/o/oauth2/v2/auth",
						},
						OAuth: config.GoogleOAuthConfig{
							ClientID:     "google-oauth-client-id",
							ClientSecret: "google-oauth-client-secret",
						},
					},
					Service: config.ServiceConfig{
						Auth: config.AuthServiceConfig{
							OAuthStateTimeout:  42 * time.Minute,
							GoogleCallbackPath: "/api/v1/google/sign-in/finish",
						},
					},
				}
			)

			for k, v := range givenEnvs {
				GinkgoT().Setenv(k, v)
			}

			cfg, err := config.LoadValidated(givenConfigDir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*cfg).Should(Equal(wantConfig))
		})

		It("loads config overwritten by local config", func() {
			var (
				givenConfigTextOverwrite = `http:
  port: 12345
service:
  auth:
    oauth-state-timeout: 1h`
			)

			var (
				wantConfig = config.Config{
					Wire: config.WireConfig{
						InitializeTimeout: 90 * time.Minute,
						ShutdownTimeout:   90 * time.Second,
					},
					Log: log.Config{
						Format: log.FormatText,
						Level:  log.LevelDebug,
						Caller: true,
					},
					HTTP: http.Config{
						Port: 12345,
					},
					Postgres: postgres.Config{
						Host:     "127.0.0.1",
						Port:     1234,
						User:     "tester",
						Password: "password",
						Database: "fake",
					},
					Redis: redis.Config{
						Addr:     "localhost:4120",
						Password: "memcached",
					},
					Google: config.GoogleConfig{
						Endpoints: config.GoogleEndpointsConfig{
							Token: "https://oauth2.googleapis.com/token",
							OAuth: "https://accounts.google.com/o/oauth2/v2/auth",
						},
						OAuth: config.GoogleOAuthConfig{
							ClientID:     "google-oauth-client-id",
							ClientSecret: "google-oauth-client-secret",
						},
					},
					Service: config.ServiceConfig{
						Auth: config.AuthServiceConfig{
							OAuthStateTimeout:  time.Hour,
							GoogleCallbackPath: "/api/v1/google/sign-in/finish",
						},
					},
				}
			)

			file := filepath.Join(givenConfigDir, "config.local.yaml")
			err := os.WriteFile(file, []byte(givenConfigTextOverwrite), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			cfg, err := config.LoadValidated(givenConfigDir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*cfg).Should(Equal(wantConfig))
		})
	})
})
