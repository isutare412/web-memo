package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type Service struct {
	kvRepository   port.KVRepository
	userRepository port.UserRepository
	googleClient   port.GoogleClient
	jwtClient      port.JWTClient

	googleOAuthEndpoint     string
	googleOAuthClientID     string
	googleOAuthCallbackPath string
	oauthStateTimeout       time.Duration
}

func NewService(
	cfg Config,
	kvRepository port.KVRepository,
	userRepository port.UserRepository,
	googleClient port.GoogleClient,
	jwtClient port.JWTClient,
) *Service {
	return &Service{
		kvRepository:   kvRepository,
		userRepository: userRepository,
		googleClient:   googleClient,

		googleOAuthEndpoint:     cfg.Google.OAuthEndpoint,
		googleOAuthClientID:     cfg.Google.OAuthClientID,
		googleOAuthCallbackPath: cfg.Google.OAuthCallbackPath,
		oauthStateTimeout:       cfg.OAuthStateTimeout,
	}
}

func (s *Service) StartGoogleSignIn(ctx context.Context, req *http.Request) (redirectURL string, err error) {
	callbackURL, err := s.getGoogleCallbackURL(req)
	if err != nil {
		return "", fmt.Errorf("getting google callback URL: %w", err)
	}

	stateID, err := s.generateOAuthStateID(ctx)
	if err != nil {
		return "", fmt.Errorf("generating oauth ID: %w", err)
	}

	oidcReq := googleOIDCRequest{
		endpoint:    s.googleOAuthEndpoint,
		clientID:    s.googleOAuthClientID,
		redirectURI: callbackURL,
		state: oauthState{
			ID:      stateID,
			Referer: req.Referer(),
		},
	}

	redirectURL, err = oidcReq.buildURL()
	if err != nil {
		return "", fmt.Errorf("building redirect URL: %w", err)
	}

	return redirectURL, nil
}

func (s *Service) FinishGoogleSignIn(ctx context.Context, req *http.Request) (redirectURL string, err error) {
	state, err := parseGoogleOAuthState(req.URL.Query())
	if err != nil {
		return "", fmt.Errorf("getting google oauth state: %w", err)
	}

	if _, err := s.kvRepository.GetThenDelete(ctx, state.ID); err != nil {
		if pkgerr.IsErrNotFound(err) {
			return "", pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "OAuth2.0 state not found",
				Origin:    err,
			}
		}
		return "", fmt.Errorf("get then deleting state: %w", err)
	}

	authCode := req.URL.Query().Get("code")
	if authCode == "" {
		return "", pkgerr.Known{
			ClientMsg: "no authorization code",
		}
	}

	callbackURL, err := s.getGoogleCallbackURL(req)
	if err != nil {
		return "", fmt.Errorf("getting google callback URL: %w", err)
	}

	tokenResp, err := s.googleClient.ExchangeAuthCode(ctx, authCode, callbackURL)
	if err != nil {
		return "", fmt.Errorf("exchanging auth code: %w", err)
	}

	idToken, err := s.jwtClient.ParseGoogleIDTokenUnverified(tokenResp.IDToken)
	if err != nil {
		return "", fmt.Errorf("parsing google ID token: %w", err)
	}

	_, err = s.userRepository.Upsert(ctx, &ent.User{
		Email:      idToken.Email,
		UserName:   idToken.Name,
		GivenName:  idToken.GivenName,
		FamilyName: idToken.FamilyName,
		PhotoURL:   idToken.PictureURL,
	})
	if err != nil {
		return "", fmt.Errorf("creating user: %w", err)
	}

	redirectURL = getBaseURL(req)
	if state.Referer != "" {
		redirectURL = state.Referer
	}

	return redirectURL, nil
}

func (s *Service) generateOAuthStateID(ctx context.Context) (string, error) {
	id := uuid.NewString()
	if err := s.kvRepository.Set(ctx, id, "", s.oauthStateTimeout); err != nil {
		return "", fmt.Errorf("setting oauth state: %w", err)
	}
	return id, nil
}

func (s *Service) getGoogleCallbackURL(r *http.Request) (string, error) {
	url, err := url.JoinPath(getBaseURL(r), s.googleOAuthCallbackPath)
	if err != nil {
		return "", err
	}
	return url, nil
}

func getBaseURL(r *http.Request) string {
	scheme := "http"
	switch {
	case r.Header.Get("X-Forwarded-Proto") == "https":
		fallthrough
	case r.TLS != nil:
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func parseGoogleOAuthState(query url.Values) (oauthState, error) {
	stateStr := query.Get(queryState)
	if stateStr == "" {
		return oauthState{}, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "state not found from query",
		}
	}

	var state oauthState
	if err := json.Unmarshal([]byte(stateStr), &state); err != nil {
		return oauthState{}, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "state is not a valid format",
			Origin:    err,
		}
	}

	return state, nil
}
