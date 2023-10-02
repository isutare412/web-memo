package config_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/config"
	"github.com/isutare412/web-memo/api/internal/log"
)

var _ = Describe("Loader", func() {
	Context("LoadValidated", func() {
		var (
			givenConfigText = `wire:
  initializeTimeout: 1h30m
  shutdownTimeout: 90s
log:
  format: text
  level: debug
  caller: true
postgres:
  host: 127.0.0.1
  port: 1234
  user: tester
  password: password
  database: fake
redis:
  addr: localhost:4120
  password: memcached`
		)

		It("loads valid config", func() {
			var (
				wantConfig = config.Config{
					Wire: config.WireConfig{
						InitializeTimeout: 90 * time.Minute,
						ShutdownTimeout:   90 * time.Second,
					},
					Log: config.LogConfig{
						Format: log.FormatText,
						Level:  log.LevelDebug,
						Caller: true,
					},
					Postgres: config.PostgresConfig{
						Host:     "127.0.0.1",
						Port:     1234,
						User:     "tester",
						Password: "password",
						Database: "fake",
					},
					Redis: config.RedisConfig{
						Addr:     "localhost:4120",
						Password: "memcached",
					},
				}
			)

			file := filepath.Join(GinkgoT().TempDir(), "config.yaml")
			err := os.WriteFile(file, []byte(givenConfigText), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			cfg, err := config.LoadValidated(file)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*cfg).Should(Equal(wantConfig))
		})

		It("loads config overwritten by env", func() {
			var (
				givenEnvs = map[string]string{
					"APP_POSTGRES_HOST":        "1.2.3.4",
					"APP_WIRE_SHUTDOWNTIMEOUT": "1s",
				}
			)

			var (
				wantConfig = config.Config{
					Wire: config.WireConfig{
						InitializeTimeout: 90 * time.Minute,
						ShutdownTimeout:   time.Second,
					},
					Log: config.LogConfig{
						Format: log.FormatText,
						Level:  log.LevelDebug,
						Caller: true,
					},
					Postgres: config.PostgresConfig{
						Host:     "1.2.3.4",
						Port:     1234,
						User:     "tester",
						Password: "password",
						Database: "fake",
					},
					Redis: config.RedisConfig{
						Addr:     "localhost:4120",
						Password: "memcached",
					},
				}
			)

			file := filepath.Join(GinkgoT().TempDir(), "config.yaml")
			err := os.WriteFile(file, []byte(givenConfigText), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			for k, v := range givenEnvs {
				GinkgoT().Setenv(k, v)
			}

			cfg, err := config.LoadValidated(file)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*cfg).Should(Equal(wantConfig))
		})
	})
})
