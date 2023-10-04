package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type Service struct {
	kvRepository port.KVRepository

	googleOAuthEndpoint     string
	googleOAuthCallbackPath string
	googleOAuthClientID     string
	googleOAuthClientSecret string
	oauthStateTimeout       time.Duration
}

func NewService(cfg Config, kvRepository port.KVRepository) *Service {
	return &Service{
		kvRepository: kvRepository,

		googleOAuthEndpoint:     cfg.Google.OAuthEndpoint,
		googleOAuthCallbackPath: cfg.Google.OAuthCallbackPath,
		googleOAuthClientID:     cfg.Google.OAuthClientID,
		googleOAuthClientSecret: cfg.Google.OAuthClientSecret,
		oauthStateTimeout:       cfg.OAuthStateTimeout,
	}
}

func (s *Service) StartGoogleSignIn(ctx context.Context, req *http.Request) (redirectURL string, err error) {
	callbackURL := &url.URL{
		Scheme: req.URL.Scheme,
		Host:   req.Host,
	}
	callbackURL = callbackURL.JoinPath(s.googleOAuthCallbackPath)

	stateID, err := s.generateOAuthStateID(ctx)
	if err != nil {
		return "", fmt.Errorf("generating oauth ID: %w", err)
	}

	oidcReq := googleOIDCRequest{
		endpoint:    s.googleOAuthEndpoint,
		clientID:    s.googleOAuthClientID,
		redirectURI: callbackURL.String(),
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

func (s *Service) generateOAuthStateID(ctx context.Context) (string, error) {
	id := uuid.NewString()
	if err := s.kvRepository.Set(ctx, id, "", s.oauthStateTimeout); err != nil {
		return "", fmt.Errorf("setting oauth state: %w", err)
	}
	return id, nil
}
