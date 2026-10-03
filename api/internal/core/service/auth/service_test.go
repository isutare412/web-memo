package auth_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port/mockport"
	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("Service", func() {
	Context("service methods", func() {
		var authService *auth.Service

		var (
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
			mockTransactionManager = mockport.NewMockTransactionManager(GinkgoT())
			mockKVRepository = mockport.NewMockKVRepository(GinkgoT())
			mockUserRepository = mockport.NewMockUserRepository(GinkgoT())
			mockGoogleClient = mockport.NewMockGoogleClient(GinkgoT())
			mockJWTClient = mockport.NewMockJWTClient(GinkgoT())

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
					Set(mock.Anything, mock.Anything, "", givenAuthConfig.OAuthStateTimeout).
					RunAndReturn(func(_ context.Context, key, _ string, _ time.Duration) error {
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

			It("stores CLI login in state value", func(ctx SpecContext) {
				var (
					givenHost        = "my-web-memo.com:1234"
					givenCallback    = "http://127.0.0.1:53682/callback"
					givenCLIState    = "abcdefghijklmnop"
					givenHTTPRequest = &http.Request{
						Host: givenHost,
						URL: &url.URL{
							Scheme: "https",
							Host:   givenHost,
							RawQuery: url.Values{
								"cliCallback": []string{givenCallback},
								"cliState":    []string{givenCLIState},
							}.Encode(),
						},
					}
					wantValue = `{"cliCallback":"http://127.0.0.1:53682/callback","cliState":"abcdefghijklmnop"}`
				)

				mockKVRepository.EXPECT().
					Set(mock.Anything, mock.Anything, wantValue, givenAuthConfig.OAuthStateTimeout).
					Return(nil)

				redirectURL, err := authService.StartGoogleSignIn(ctx, givenHTTPRequest)
				Expect(err).ShouldNot(HaveOccurred())

				unescapedURL, err := url.QueryUnescape(redirectURL)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(unescapedURL).ShouldNot(ContainSubstring("cliCallback"))
				Expect(unescapedURL).ShouldNot(ContainSubstring(givenCallback))
				Expect(unescapedURL).ShouldNot(ContainSubstring(givenCLIState))
			})

			It("accepts localhost CLI callback", func(ctx SpecContext) {
				var (
					givenHost        = "my-web-memo.com:1234"
					givenHTTPRequest = &http.Request{
						Host: givenHost,
						URL: &url.URL{
							Scheme: "https",
							Host:   givenHost,
							RawQuery: url.Values{
								"cliCallback": []string{"http://localhost:53682/callback"},
								"cliState":    []string{"abcdefghijklmnop"},
							}.Encode(),
						},
					}
					wantValue = `{"cliCallback":"http://localhost:53682/callback","cliState":"abcdefghijklmnop"}`
				)

				mockKVRepository.EXPECT().
					Set(mock.Anything, mock.Anything, wantValue, givenAuthConfig.OAuthStateTimeout).
					Return(nil)

				_, err := authService.StartGoogleSignIn(ctx, givenHTTPRequest)
				Expect(err).ShouldNot(HaveOccurred())
			})

			DescribeTable("rejects invalid CLI login parameters",
				func(ctx SpecContext, callback, cliState string) {
					query := url.Values{}
					if callback != "" {
						query.Set("cliCallback", callback)
					}
					if cliState != "" {
						query.Set("cliState", cliState)
					}
					givenHTTPRequest := &http.Request{
						Host: "my-web-memo.com",
						URL: &url.URL{
							Scheme:   "https",
							Host:     "my-web-memo.com",
							RawQuery: query.Encode(),
						},
					}

					_, err := authService.StartGoogleSignIn(ctx, givenHTTPRequest)
					Expect(err).Should(HaveOccurred())

					var known pkgerr.Known
					Expect(errors.As(err, &known)).Should(BeTrue())
					Expect(known.Code).Should(Equal(pkgerr.CodeBadRequest))
					mockKVRepository.AssertNotCalled(GinkgoT(), "Set",
						mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				},
				Entry("https scheme", "https://127.0.0.1:1/callback", "abcdefghijklmnop"),
				Entry("non-loopback host", "http://evil.com:1/callback", "abcdefghijklmnop"),
				Entry("missing port", "http://127.0.0.1/callback", "abcdefghijklmnop"),
				Entry("wrong path", "http://127.0.0.1:1/other", "abcdefghijklmnop"),
				Entry("query string", "http://127.0.0.1:1/callback?x=1", "abcdefghijklmnop"),
				Entry("fragment", "http://127.0.0.1:1/callback#f", "abcdefghijklmnop"),
				Entry("missing cliState", "http://127.0.0.1:1/callback", ""),
				Entry("short cliState", "http://127.0.0.1:1/callback", "abcdefghijklmno"),
				Entry("cliState with dot", "http://127.0.0.1:1/callback", "abcdefghijklmno."),
				Entry("missing cliCallback", "", "abcdefghijklmnop"),
				Entry("localhost prefix host", "http://localhost.evil.com:1/callback", "abcdefghijklmnop"),
				Entry("userinfo host confusion", "http://127.0.0.1:1@evil.com/callback", "abcdefghijklmnop"),
				Entry("ipv6 loopback", "http://[::1]:1/callback", "abcdefghijklmnop"),
				Entry("localhost trailing dot", "http://localhost.:1/callback", "abcdefghijklmnop"),
			)
		})

		Context("FinishGoogleSignIn", func() {
			It("signs app ID token using google tokens", func(ctx SpecContext) {
				var (
					givenHost     = "localhost:42"
					givenStateID  = uuid.NewString()
					givenReferer  = "http://localhost:1234/foo/page"
					givenState    = fmt.Sprintf(`{"id":"%s","referer":"%s"}`, givenStateID, givenReferer)
					givenAuthCode = "auth-code-from-google"
					givenURLQuery = url.Values{
						"state": []string{givenState},
						"code":  []string{givenAuthCode},
					}
					givenHTTPRequest = &http.Request{
						Host: givenHost,
						URL: &url.URL{
							RawQuery: givenURLQuery.Encode(),
						},
					}
					givenGoogleIDToken = "id-token-from-google"
					givenUser          = &ent.User{
						ID:         uuid.New(),
						Email:      "foo@gmail.com",
						UserName:   "Alice Bob",
						GivenName:  "Alice",
						FamilyName: "Bob",
						PhotoURL:   "https://my-pic.com/foo",
						Type:       enum.UserTypeClient,
					}
					givenAppIDToken = "app-id-token"
				)

				mockKVRepository.EXPECT().
					GetThenDelete(mock.Anything, givenStateID).
					Return("", nil)

				mockGoogleClient.EXPECT().
					ExchangeAuthCode(mock.Anything, givenAuthCode, mock.Anything).
					RunAndReturn(func(_ context.Context, _, redirectURI string) (model.GoogleTokenResponse, error) {
						baseURL := fmt.Sprintf("http://%s", givenHost)
						callbackURL, err := url.JoinPath(baseURL, givenAuthConfig.Google.OAuthCallbackPath)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(redirectURI).Should(Equal(callbackURL))

						return model.GoogleTokenResponse{
							IDToken: givenGoogleIDToken,
						}, nil
					})

				mockJWTClient.EXPECT().
					ParseGoogleIDTokenUnverified(givenGoogleIDToken).
					Return(&model.GoogleIDToken{
						Email:      givenUser.Email,
						Name:       givenUser.UserName,
						GivenName:  givenUser.GivenName,
						FamilyName: givenUser.FamilyName,
						PictureURL: givenUser.PhotoURL,
					}, nil)

				mockTransactionManager.EXPECT().
					WithTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					})

				mockUserRepository.EXPECT().
					FindByEmail(mock.Anything, givenUser.Email).
					Return(nil, pkgerr.Known{Code: pkgerr.CodeNotFound})

				mockUserRepository.EXPECT().
					Upsert(mock.Anything, mock.Anything).
					Return(givenUser, nil)

				mockJWTClient.EXPECT().
					SignAppIDToken(mock.Anything).
					RunAndReturn(func(t *model.AppIDToken) (token *model.AppIDToken, tokenString string, err error) {
						Expect(t.UserID).Should(Equal(givenUser.ID))
						Expect(t.Email).Should(Equal(givenUser.Email))
						Expect(t.UserName).Should(Equal(givenUser.UserName))
						Expect(t.GivenName).Should(Equal(givenUser.GivenName))
						Expect(t.FamilyName).Should(Equal(givenUser.FamilyName))
						Expect(t.PhotoURL).Should(Equal(givenUser.PhotoURL))
						return nil, givenAppIDToken, nil
					})

				redirectURL, appIDToken, err := authService.FinishGoogleSignIn(ctx, givenHTTPRequest)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(redirectURL).Should(Equal(givenReferer))
				Expect(appIDToken).Should(Equal(givenAppIDToken))
			})
		})
	})
})
