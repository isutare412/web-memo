package auth_test

import (
	"context"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/isutare412/web-memo/api/internal/core/port/mockport"
	"github.com/isutare412/web-memo/api/internal/core/service/auth"
)

var _ = Describe("Service", func() {
	Context("service methods", func() {
		var authService *auth.Service

		var (
			mockController         *gomock.Controller
			mockTransactionManager *mockport.MockTransactionManager
			mockKVRepository       *mockport.MockKVRepository
			mockUserRepository     *mockport.MockUserRepository
			mockGoogleClient       *mockport.MockGoogleClient
			mockJWTClient          *mockport.MockJWTClient
		)

		var (
			givenAuthConfig = auth.Config{
				Google: auth.GoogleConfig{
					OAuthEndpoint:     "https://accounts.google.com/o/oauth2/v2/auth",
					OAuthClientID:     "google-client-id",
					OAuthCallbackPath: "/google/callback",
				},
				OAuthStateTimeout: time.Second,
			}
		)

		BeforeEach(func() {
			mockController = gomock.NewController(GinkgoT())
			mockTransactionManager = mockport.NewMockTransactionManager(mockController)
			mockKVRepository = mockport.NewMockKVRepository(mockController)
			mockUserRepository = mockport.NewMockUserRepository(mockController)
			mockGoogleClient = mockport.NewMockGoogleClient(mockController)
			mockJWTClient = mockport.NewMockJWTClient(mockController)

			authService = auth.NewService(
				givenAuthConfig, mockTransactionManager, mockKVRepository, mockUserRepository,
				mockGoogleClient, mockJWTClient)
		})

		Context("StartGoogleSignIn", func() {
			It("builds redirect URL as expected", func(ctx SpecContext) {
				var (
					givenHost        = "my-web-memo.com:1234"
					givenReferer     = "https://my-web-app"
					givenHTTPRequest = &http.Request{
						Host: givenHost,
						URL: &url.URL{
							Scheme: "https",
							Host:   givenHost,
						},
						Header: http.Header{
							"Referer": []string{givenReferer},
						},
					}
				)

				var (
					gotStateID string
				)

				mockKVRepository.EXPECT().
					Set(gomock.Any(), gomock.Any(), "", givenAuthConfig.OAuthStateTimeout).
					DoAndReturn(func(_ context.Context, key, _ string, _ time.Duration) error {
						gotStateID = key
						return nil
					})

				redirectURL, err := authService.StartGoogleSignIn(ctx, givenHTTPRequest)
				Expect(err).ShouldNot(HaveOccurred())

				unescapedURL, err := url.QueryUnescape(redirectURL)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(unescapedURL).Should(ContainSubstring(gotStateID))
				Expect(unescapedURL).Should(ContainSubstring(givenHost))
				Expect(unescapedURL).Should(ContainSubstring(givenReferer))
				Expect(unescapedURL).Should(ContainSubstring(givenAuthConfig.Google.OAuthClientID))
				Expect(unescapedURL).Should(ContainSubstring(givenAuthConfig.Google.OAuthCallbackPath))
				Expect(unescapedURL).Should(ContainSubstring(givenAuthConfig.Google.OAuthEndpoint))
			})
		})
	})
})
