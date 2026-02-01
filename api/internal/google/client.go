package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/trace"
)

type Client struct {
	httpClient *http.Client

	tokenEndpoint     string
	oauthClientID     string
	oauthClientSecret string
}

func NewClient(cfg ClientConfig) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 100

	return &Client{
		httpClient: &http.Client{Transport: transport},

		tokenEndpoint:     cfg.TokenEndpoint,
		oauthClientID:     cfg.OAuthClientID,
		oauthClientSecret: cfg.OAuthClientSecret,
	}
}

func (c *Client) ExchangeAuthCode(ctx context.Context, code, redirectURI string) (model.GoogleTokenResponse, error) {
	ctx, span := trace.StartSpan(ctx, "google.Client.ExchangeAuthCode")
	defer span.End()

	body := url.Values{}
	body.Add("client_id", c.oauthClientID)
	body.Add("client_secret", c.oauthClientSecret)
	body.Add("grant_type", "authorization_code")
	body.Add("code", code)
	body.Add("redirect_uri", redirectURI)
	encodedBody := body.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", c.tokenEndpoint, strings.NewReader(encodedBody))
	if err != nil {
		return model.GoogleTokenResponse{}, fmt.Errorf("building request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rawResp, err := c.httpClient.Do(req)
	switch {
	case err != nil:
		return model.GoogleTokenResponse{}, fmt.Errorf("doing HTTP request: %w", err)
	case rawResp.StatusCode >= 400:
		return model.GoogleTokenResponse{}, fmt.Errorf("got error from Google token API; statusCode(%d)", rawResp.StatusCode)
	}
	defer rawResp.Body.Close()

	var resp googleOAuthTokens
	if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
		return model.GoogleTokenResponse{}, fmt.Errorf("decoding auth code response: %w", err)
	}

	return resp.ToTokenResponse(), nil
}
