package auth

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

// googleOIDCRequest holds every data to start OpenID Connect flow.
// Check this out https://developers.google.com/identity/openid-connect/openid-connect#authenticationuriparameters.
type googleOIDCRequest struct {
	endpoint    string
	clientID    string
	redirectURI string
	state       oauthState
}

func (r *googleOIDCRequest) buildURL() (string, error) {
	uri, err := url.Parse(r.endpoint)
	if err != nil {
		return "", fmt.Errorf("parsing Google OIDC endpoint: %w", err)
	}

	stateBytes, err := json.Marshal(&r.state)
	if err != nil {
		return "", fmt.Errorf("marshaling oauth state: %w", err)
	}

	q := url.Values{}
	q.Add("access_type", "online")
	q.Add("response_type", "code")
	q.Add("scope", "openid profile email")
	q.Add("prompt", "none")
	q.Add("nonce", uuid.NewString())
	q.Add("client_id", r.clientID)
	q.Add("redirect_uri", r.redirectURI)
	q.Add("state", string(stateBytes))
	uri.RawQuery = q.Encode()

	return uri.String(), nil
}

type oauthState struct {
	ID      string `json:"id"`
	Referer string `json:"referer,omitempty"`
}
