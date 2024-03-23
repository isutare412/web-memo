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
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type Service struct {
	transactionManager port.TransactionManager
	kvRepository       port.KVRepository
	userRepository     port.UserRepository
	googleClient       port.GoogleClient
	jwtClient          port.JWTClient

	googleOAuthEndpoint     string
	googleOAuthClientID     string
	googleOAuthCallbackPath string
	oauthStateTimeout       time.Duration
}

func NewService(
	cfg Config,
	transactionManager port.TransactionManager,
	kvRepository port.KVRepository,
	userRepository port.UserRepository,
	googleClient port.GoogleClient,
	jwtClient port.JWTClient,
) *Service {
	return &Service{
		transactionManager: transactionManager,
		kvRepository:       kvRepository,
		userRepository:     userRepository,
		googleClient:       googleClient,
		jwtClient:          jwtClient,

		googleOAuthEndpoint:     cfg.Google.OAuthEndpoint,
		googleOAuthClientID:     cfg.Google.OAuthClientID,
		googleOAuthCallbackPath: cfg.Google.OAuthCallbackPath,
		oauthStateTimeout:       cfg.OAuthStateTimeout,
	}
}

func (s *Service) VerifyAppIDToken(tokenString string) (*model.AppIDToken, error) {
	return s.jwtClient.VerifyAppIDToken(tokenString)
}

func (s *Service) RefreshAppIDToken(
	ctx context.Context,
	tokenString string,
) (newToken *model.AppIDToken, newTokenString string, err error) {
	token, err := s.jwtClient.VerifyAppIDToken(tokenString)
	if err != nil {
		return nil, "", fmt.Errorf("verifying app id token: %w", err)
	}

	user, err := s.userRepository.FindByEmail(ctx, token.Email)
	if err != nil {
		return nil, "", fmt.Errorf("finding user: %w", err)
	}

	newToken, newTokenString, err = s.jwtClient.SignAppIDToken(model.NewAppIDToken(user))
	if err != nil {
		return nil, "", fmt.Errorf("signing app id token: %w", err)
	}

	return newToken, newTokenString, nil
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

func (s *Service) FinishGoogleSignIn(
	ctx context.Context,
	req *http.Request,
) (redirectURL, appTokenString string, err error) {
	state, err := parseGoogleOAuthState(req.URL.Query())
	if err != nil {
		return "", "", fmt.Errorf("getting google oauth state: %w", err)
	}

	if _, err := s.kvRepository.GetThenDelete(ctx, state.ID); err != nil {
		if pkgerr.IsErrNotFound(err) {
			return "", "", pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "OAuth2.0 state not found",
				Origin:    err,
			}
		}
		return "", "", fmt.Errorf("get then deleting state: %w", err)
	}

	authCode := req.URL.Query().Get("code")
	if authCode == "" {
		return "", "", pkgerr.Known{
			ClientMsg: "no authorization code",
		}
	}

	callbackURL, err := s.getGoogleCallbackURL(req)
	if err != nil {
		return "", "", fmt.Errorf("getting google callback URL: %w", err)
	}

	tokenResp, err := s.googleClient.ExchangeAuthCode(ctx, authCode, callbackURL)
	if err != nil {
		return "", "", fmt.Errorf("exchanging auth code: %w", err)
	}

	idToken, err := s.jwtClient.ParseGoogleIDTokenUnverified(tokenResp.IDToken)
	if err != nil {
		return "", "", fmt.Errorf("parsing google ID token: %w", err)
	}

	var userUpserted *ent.User
	err = s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		userToCreate := &ent.User{
			Email:      idToken.Email,
			UserName:   idToken.Name,
			GivenName:  idToken.GivenName,
			FamilyName: idToken.FamilyName,
			PhotoURL:   idToken.PictureURL,
			Type:       enum.UserTypeClient,
		}

		userFound, err := s.userRepository.FindByEmail(ctx, idToken.Email)
		switch {
		case err != nil && !pkgerr.IsErrNotFound(err):
			return fmt.Errorf("finding user by email: %w", err)
		case err == nil:
			userToCreate.Type = userFound.Type
		}

		user, err := s.userRepository.Upsert(ctx, userToCreate)
		if err != nil {
			return fmt.Errorf("upserting user: %w", err)
		}

		userUpserted = user
		return nil
	})
	if err != nil {
		return "", "", fmt.Errorf("during transaction: %w", err)
	}

	_, appTokenString, err = s.jwtClient.SignAppIDToken(model.NewAppIDToken(userUpserted))
	if err != nil {
		return "", "", fmt.Errorf("signing app ID token: %w", err)
	}

	redirectURL = getBaseURL(req)
	if state.Referer != "" {
		redirectURL = state.Referer
	}

	return redirectURL, appTokenString, nil
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
