package config_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/config"
	"github.com/isutare412/web-memo/api/internal/jwt"
	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
	"github.com/isutare412/web-memo/api/internal/tracing"
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
trace:
  enabled: false
  service-name: webmemo-api-test
  sampling-ratio: 0.5
  otlp-grpc-endpoint: localhost:4317
web:
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
jwt:
  expiration: 720h
  active-key-pair: test
  key-pairs:
    - name: 'test'
      private: |
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpQIBAAKCAQEAv+oq9v//ugc1cjhRlxJgd+R8wo72XtlfI5VWkKZNq/bnpsCo
        bOE11viKPrakQS9fwh8WKvesOz63y/RIMGPmRVydo9kZoPTlefTjENiNSHDH0T+Y
        L8JAscbczHgNvfqfy0UM+eC21/T435zShNKcUUyyGR6s99fhZVIzlby3H1A/cQoV
        TiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d9ycjoO74QP5RQdwLU87b3TXLSopa
        mCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+o0oHxOwz56neM+plxYKvLKmUsdyE
        Ul2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4zQIDAQABAoIBAQCGH7fLS/qDHoqh
        uu56sGMvJ0ZyCsvwWeZ9zd7j1PYvmq0nAzoybercxHKJhcehruQznNo3SUTbWufE
        6IKTHx5Nl36shgu9S6oc46LVoSKMYBWmDdXketQP6rVhSP4BqeiHfUimUgA3SYOt
        c8JFBZQt1XYazC+CPyPNVferTGqGvK7X8jZkPCQyQxj+EAZ7k0cvZLz8ifCyf/E9
        qQh1sl8qdgAZuf+1KQUNEND9soD0R4RkbftiJjUltI1mFmuiEIb0MF88Of1MGcF1
        2Aku5CVyr8FQRhvbRSS8ZPDnOEnjfFTBF4P/bP5ToPdEpx7G3PMgm6pCQ0rXgFUe
        F500fxcBAoGBAOItTIye/dycARhZUyWwv8pG5gBqqddIehjmLqo7euaWws1YcJL9
        q8MXKIR51c75SjpC4ya9qSYHzc9SBP77S++l7ROT7C9hfwRDtNuNxM1FDyOiZq7F
        0yJqqKDQ1VRf/Zen5Jg0raQtKMND/m1S+dA/qgGoZV6MaE+XM2AATvYhAoGBANk4
        VRRkLtoVNl3vRx82EBTA40nPqX4P9j/80isXlTkD3ynwfASHS5FL24uwW6Fe4bEu
        ATDhEiDnNROMHSq3H0OfYpO6YJ70zzfwZ1aKu6CU5cnxfuPzvnnea3JplkPIvFo8
        3OXOViIYoz34rU0ilWEFjuFovY4p4ob907amOdUtAoGBALpMkccqlvSGU6iYuxJK
        mk+lQoKJWUiI3Hlx8HIr+DnDaMX32RJafIZ/ptIoAOMxF+ERg0U/5/n5Z58jchYN
        LClDxRnhOCR27Ea49ln6Vma2QZgahvXi4Nxyel+sZGvRfXLTyklM6tJWmELu2L14
        IWlVZ1ViPc05Xhpg8uJanq/BAoGBAIfo3CrXCA2BijO528kmfWdOzKdJHCZ4/D1L
        BYDaz44N4xqNkjsPH/P3/5TmMl7ES/gc7bfUixA1OZtSZolsbE5WMkp2KbArQmAg
        tbeLNBwkLaZtyFP+FOaRiK7ca51bwqW/QQM0V+Ybfj/vERebFNXQsXZNn5SMlmSZ
        +lZkqPi9AoGAOEiNa31aeFN+h9mMwIGpEZLXQe+EJXPKPL4Svaa4UMviGgxfpHss
        acuOiZsNxxGDhzRrz4YhqfCTcnoc4gL3TVwQN38fzOTVkQtD73ENAI+fUyZ32gLK
        EppZW3NlC2IZVtJMB9jAk1WlsppD0plFZDmy4IB8WbfVARmo0DWhTE4=
        -----END RSA PRIVATE KEY-----
      public: |
        -----BEGIN PUBLIC KEY-----
        MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv+oq9v//ugc1cjhRlxJg
        d+R8wo72XtlfI5VWkKZNq/bnpsCobOE11viKPrakQS9fwh8WKvesOz63y/RIMGPm
        RVydo9kZoPTlefTjENiNSHDH0T+YL8JAscbczHgNvfqfy0UM+eC21/T435zShNKc
        UUyyGR6s99fhZVIzlby3H1A/cQoVTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d
        9ycjoO74QP5RQdwLU87b3TXLSopamCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+
        o0oHxOwz56neM+plxYKvLKmUsdyEUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4
        zQIDAQAB
        -----END PUBLIC KEY-----
oauth:
  state-timeout: 42m
  callback-path: /api/v1/google/sign-in/finish
cron:
  tag-cleanup-interval: 19m
imageer:
  base-url: https://imageer.example.com
  api-key: test-api-key
  project-id: 12345678-1234-1234-1234-123456789012`
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
					Trace: tracing.Config{
						ServiceName:      "webmemo-api-test",
						SamplingRatio:    0.5,
						OTLPGRPCEndpoint: "localhost:4317",
					},
					Web: config.WebConfig{
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
					JWT: jwt.Config{
						ActiveKeyPair: "test",
						KeyPairs: []jwt.RSAKeyBytesPair{{
							Name:    "test",
							Private: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAv+oq9v//ugc1cjhRlxJgd+R8wo72XtlfI5VWkKZNq/bnpsCo\nbOE11viKPrakQS9fwh8WKvesOz63y/RIMGPmRVydo9kZoPTlefTjENiNSHDH0T+Y\nL8JAscbczHgNvfqfy0UM+eC21/T435zShNKcUUyyGR6s99fhZVIzlby3H1A/cQoV\nTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d9ycjoO74QP5RQdwLU87b3TXLSopa\nmCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+o0oHxOwz56neM+plxYKvLKmUsdyE\nUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4zQIDAQABAoIBAQCGH7fLS/qDHoqh\nuu56sGMvJ0ZyCsvwWeZ9zd7j1PYvmq0nAzoybercxHKJhcehruQznNo3SUTbWufE\n6IKTHx5Nl36shgu9S6oc46LVoSKMYBWmDdXketQP6rVhSP4BqeiHfUimUgA3SYOt\nc8JFBZQt1XYazC+CPyPNVferTGqGvK7X8jZkPCQyQxj+EAZ7k0cvZLz8ifCyf/E9\nqQh1sl8qdgAZuf+1KQUNEND9soD0R4RkbftiJjUltI1mFmuiEIb0MF88Of1MGcF1\n2Aku5CVyr8FQRhvbRSS8ZPDnOEnjfFTBF4P/bP5ToPdEpx7G3PMgm6pCQ0rXgFUe\nF500fxcBAoGBAOItTIye/dycARhZUyWwv8pG5gBqqddIehjmLqo7euaWws1YcJL9\nq8MXKIR51c75SjpC4ya9qSYHzc9SBP77S++l7ROT7C9hfwRDtNuNxM1FDyOiZq7F\n0yJqqKDQ1VRf/Zen5Jg0raQtKMND/m1S+dA/qgGoZV6MaE+XM2AATvYhAoGBANk4\nVRRkLtoVNl3vRx82EBTA40nPqX4P9j/80isXlTkD3ynwfASHS5FL24uwW6Fe4bEu\nATDhEiDnNROMHSq3H0OfYpO6YJ70zzfwZ1aKu6CU5cnxfuPzvnnea3JplkPIvFo8\n3OXOViIYoz34rU0ilWEFjuFovY4p4ob907amOdUtAoGBALpMkccqlvSGU6iYuxJK\nmk+lQoKJWUiI3Hlx8HIr+DnDaMX32RJafIZ/ptIoAOMxF+ERg0U/5/n5Z58jchYN\nLClDxRnhOCR27Ea49ln6Vma2QZgahvXi4Nxyel+sZGvRfXLTyklM6tJWmELu2L14\nIWlVZ1ViPc05Xhpg8uJanq/BAoGBAIfo3CrXCA2BijO528kmfWdOzKdJHCZ4/D1L\nBYDaz44N4xqNkjsPH/P3/5TmMl7ES/gc7bfUixA1OZtSZolsbE5WMkp2KbArQmAg\ntbeLNBwkLaZtyFP+FOaRiK7ca51bwqW/QQM0V+Ybfj/vERebFNXQsXZNn5SMlmSZ\n+lZkqPi9AoGAOEiNa31aeFN+h9mMwIGpEZLXQe+EJXPKPL4Svaa4UMviGgxfpHss\nacuOiZsNxxGDhzRrz4YhqfCTcnoc4gL3TVwQN38fzOTVkQtD73ENAI+fUyZ32gLK\nEppZW3NlC2IZVtJMB9jAk1WlsppD0plFZDmy4IB8WbfVARmo0DWhTE4=\n-----END RSA PRIVATE KEY-----\n",
							Public:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv+oq9v//ugc1cjhRlxJg\nd+R8wo72XtlfI5VWkKZNq/bnpsCobOE11viKPrakQS9fwh8WKvesOz63y/RIMGPm\nRVydo9kZoPTlefTjENiNSHDH0T+YL8JAscbczHgNvfqfy0UM+eC21/T435zShNKc\nUUyyGR6s99fhZVIzlby3H1A/cQoVTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d\n9ycjoO74QP5RQdwLU87b3TXLSopamCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+\no0oHxOwz56neM+plxYKvLKmUsdyEUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4\nzQIDAQAB\n-----END PUBLIC KEY-----\n",
						}},
						Expiration: 720 * time.Hour,
					},
					OAuth: config.OAuthConfig{
						StateTimeout: 42 * time.Minute,
						CallbackPath: "/api/v1/google/sign-in/finish",
					},
					Cron: config.CronConfig{
						TagCleanupInterval: 19 * time.Minute,
					},
					Imageer: config.ImageerConfig{
						BaseURL:   "https://imageer.example.com",
						APIKey:    "test-api-key",
						ProjectID: "12345678-1234-1234-1234-123456789012",
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
					Trace: tracing.Config{
						ServiceName:      "webmemo-api-test",
						SamplingRatio:    0.5,
						OTLPGRPCEndpoint: "localhost:4317",
					},
					Web: config.WebConfig{
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
					JWT: jwt.Config{
						ActiveKeyPair: "test",
						KeyPairs: []jwt.RSAKeyBytesPair{{
							Name:    "test",
							Private: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAv+oq9v//ugc1cjhRlxJgd+R8wo72XtlfI5VWkKZNq/bnpsCo\nbOE11viKPrakQS9fwh8WKvesOz63y/RIMGPmRVydo9kZoPTlefTjENiNSHDH0T+Y\nL8JAscbczHgNvfqfy0UM+eC21/T435zShNKcUUyyGR6s99fhZVIzlby3H1A/cQoV\nTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d9ycjoO74QP5RQdwLU87b3TXLSopa\nmCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+o0oHxOwz56neM+plxYKvLKmUsdyE\nUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4zQIDAQABAoIBAQCGH7fLS/qDHoqh\nuu56sGMvJ0ZyCsvwWeZ9zd7j1PYvmq0nAzoybercxHKJhcehruQznNo3SUTbWufE\n6IKTHx5Nl36shgu9S6oc46LVoSKMYBWmDdXketQP6rVhSP4BqeiHfUimUgA3SYOt\nc8JFBZQt1XYazC+CPyPNVferTGqGvK7X8jZkPCQyQxj+EAZ7k0cvZLz8ifCyf/E9\nqQh1sl8qdgAZuf+1KQUNEND9soD0R4RkbftiJjUltI1mFmuiEIb0MF88Of1MGcF1\n2Aku5CVyr8FQRhvbRSS8ZPDnOEnjfFTBF4P/bP5ToPdEpx7G3PMgm6pCQ0rXgFUe\nF500fxcBAoGBAOItTIye/dycARhZUyWwv8pG5gBqqddIehjmLqo7euaWws1YcJL9\nq8MXKIR51c75SjpC4ya9qSYHzc9SBP77S++l7ROT7C9hfwRDtNuNxM1FDyOiZq7F\n0yJqqKDQ1VRf/Zen5Jg0raQtKMND/m1S+dA/qgGoZV6MaE+XM2AATvYhAoGBANk4\nVRRkLtoVNl3vRx82EBTA40nPqX4P9j/80isXlTkD3ynwfASHS5FL24uwW6Fe4bEu\nATDhEiDnNROMHSq3H0OfYpO6YJ70zzfwZ1aKu6CU5cnxfuPzvnnea3JplkPIvFo8\n3OXOViIYoz34rU0ilWEFjuFovY4p4ob907amOdUtAoGBALpMkccqlvSGU6iYuxJK\nmk+lQoKJWUiI3Hlx8HIr+DnDaMX32RJafIZ/ptIoAOMxF+ERg0U/5/n5Z58jchYN\nLClDxRnhOCR27Ea49ln6Vma2QZgahvXi4Nxyel+sZGvRfXLTyklM6tJWmELu2L14\nIWlVZ1ViPc05Xhpg8uJanq/BAoGBAIfo3CrXCA2BijO528kmfWdOzKdJHCZ4/D1L\nBYDaz44N4xqNkjsPH/P3/5TmMl7ES/gc7bfUixA1OZtSZolsbE5WMkp2KbArQmAg\ntbeLNBwkLaZtyFP+FOaRiK7ca51bwqW/QQM0V+Ybfj/vERebFNXQsXZNn5SMlmSZ\n+lZkqPi9AoGAOEiNa31aeFN+h9mMwIGpEZLXQe+EJXPKPL4Svaa4UMviGgxfpHss\nacuOiZsNxxGDhzRrz4YhqfCTcnoc4gL3TVwQN38fzOTVkQtD73ENAI+fUyZ32gLK\nEppZW3NlC2IZVtJMB9jAk1WlsppD0plFZDmy4IB8WbfVARmo0DWhTE4=\n-----END RSA PRIVATE KEY-----\n",
							Public:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv+oq9v//ugc1cjhRlxJg\nd+R8wo72XtlfI5VWkKZNq/bnpsCobOE11viKPrakQS9fwh8WKvesOz63y/RIMGPm\nRVydo9kZoPTlefTjENiNSHDH0T+YL8JAscbczHgNvfqfy0UM+eC21/T435zShNKc\nUUyyGR6s99fhZVIzlby3H1A/cQoVTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d\n9ycjoO74QP5RQdwLU87b3TXLSopamCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+\no0oHxOwz56neM+plxYKvLKmUsdyEUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4\nzQIDAQAB\n-----END PUBLIC KEY-----\n",
						}},
						Expiration: 720 * time.Hour,
					},
					OAuth: config.OAuthConfig{
						StateTimeout: 42 * time.Minute,
						CallbackPath: "/api/v1/google/sign-in/finish",
					},
					Cron: config.CronConfig{
						TagCleanupInterval: 19 * time.Minute,
					},
					Imageer: config.ImageerConfig{
						BaseURL:   "https://imageer.example.com",
						APIKey:    "test-api-key",
						ProjectID: "12345678-1234-1234-1234-123456789012",
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
				givenConfigTextOverwrite = `web:
  port: 12345
oauth:
  state-timeout: 1h`
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
					Trace: tracing.Config{
						ServiceName:      "webmemo-api-test",
						SamplingRatio:    0.5,
						OTLPGRPCEndpoint: "localhost:4317",
					},
					Web: config.WebConfig{
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
					JWT: jwt.Config{
						ActiveKeyPair: "test",
						KeyPairs: []jwt.RSAKeyBytesPair{{
							Name:    "test",
							Private: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAv+oq9v//ugc1cjhRlxJgd+R8wo72XtlfI5VWkKZNq/bnpsCo\nbOE11viKPrakQS9fwh8WKvesOz63y/RIMGPmRVydo9kZoPTlefTjENiNSHDH0T+Y\nL8JAscbczHgNvfqfy0UM+eC21/T435zShNKcUUyyGR6s99fhZVIzlby3H1A/cQoV\nTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d9ycjoO74QP5RQdwLU87b3TXLSopa\nmCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+o0oHxOwz56neM+plxYKvLKmUsdyE\nUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4zQIDAQABAoIBAQCGH7fLS/qDHoqh\nuu56sGMvJ0ZyCsvwWeZ9zd7j1PYvmq0nAzoybercxHKJhcehruQznNo3SUTbWufE\n6IKTHx5Nl36shgu9S6oc46LVoSKMYBWmDdXketQP6rVhSP4BqeiHfUimUgA3SYOt\nc8JFBZQt1XYazC+CPyPNVferTGqGvK7X8jZkPCQyQxj+EAZ7k0cvZLz8ifCyf/E9\nqQh1sl8qdgAZuf+1KQUNEND9soD0R4RkbftiJjUltI1mFmuiEIb0MF88Of1MGcF1\n2Aku5CVyr8FQRhvbRSS8ZPDnOEnjfFTBF4P/bP5ToPdEpx7G3PMgm6pCQ0rXgFUe\nF500fxcBAoGBAOItTIye/dycARhZUyWwv8pG5gBqqddIehjmLqo7euaWws1YcJL9\nq8MXKIR51c75SjpC4ya9qSYHzc9SBP77S++l7ROT7C9hfwRDtNuNxM1FDyOiZq7F\n0yJqqKDQ1VRf/Zen5Jg0raQtKMND/m1S+dA/qgGoZV6MaE+XM2AATvYhAoGBANk4\nVRRkLtoVNl3vRx82EBTA40nPqX4P9j/80isXlTkD3ynwfASHS5FL24uwW6Fe4bEu\nATDhEiDnNROMHSq3H0OfYpO6YJ70zzfwZ1aKu6CU5cnxfuPzvnnea3JplkPIvFo8\n3OXOViIYoz34rU0ilWEFjuFovY4p4ob907amOdUtAoGBALpMkccqlvSGU6iYuxJK\nmk+lQoKJWUiI3Hlx8HIr+DnDaMX32RJafIZ/ptIoAOMxF+ERg0U/5/n5Z58jchYN\nLClDxRnhOCR27Ea49ln6Vma2QZgahvXi4Nxyel+sZGvRfXLTyklM6tJWmELu2L14\nIWlVZ1ViPc05Xhpg8uJanq/BAoGBAIfo3CrXCA2BijO528kmfWdOzKdJHCZ4/D1L\nBYDaz44N4xqNkjsPH/P3/5TmMl7ES/gc7bfUixA1OZtSZolsbE5WMkp2KbArQmAg\ntbeLNBwkLaZtyFP+FOaRiK7ca51bwqW/QQM0V+Ybfj/vERebFNXQsXZNn5SMlmSZ\n+lZkqPi9AoGAOEiNa31aeFN+h9mMwIGpEZLXQe+EJXPKPL4Svaa4UMviGgxfpHss\nacuOiZsNxxGDhzRrz4YhqfCTcnoc4gL3TVwQN38fzOTVkQtD73ENAI+fUyZ32gLK\nEppZW3NlC2IZVtJMB9jAk1WlsppD0plFZDmy4IB8WbfVARmo0DWhTE4=\n-----END RSA PRIVATE KEY-----\n",
							Public:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv+oq9v//ugc1cjhRlxJg\nd+R8wo72XtlfI5VWkKZNq/bnpsCobOE11viKPrakQS9fwh8WKvesOz63y/RIMGPm\nRVydo9kZoPTlefTjENiNSHDH0T+YL8JAscbczHgNvfqfy0UM+eC21/T435zShNKc\nUUyyGR6s99fhZVIzlby3H1A/cQoVTiyuuWU5DdH/TT1ejx/kbpFtC+RHTPddoO0d\n9ycjoO74QP5RQdwLU87b3TXLSopamCJIk47QddBvV491QCSUURzY7HkUgZu0OB8+\no0oHxOwz56neM+plxYKvLKmUsdyEUl2anNn7qCHkxRq5HBQR1myWGAtNg1vF9AO4\nzQIDAQAB\n-----END PUBLIC KEY-----\n",
						}},
						Expiration: 720 * time.Hour,
					},
					OAuth: config.OAuthConfig{
						StateTimeout: time.Hour,
						CallbackPath: "/api/v1/google/sign-in/finish",
					},
					Cron: config.CronConfig{
						TagCleanupInterval: 19 * time.Minute,
					},
					Imageer: config.ImageerConfig{
						BaseURL:   "https://imageer.example.com",
						APIKey:    "test-api-key",
						ProjectID: "12345678-1234-1234-1234-123456789012",
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
