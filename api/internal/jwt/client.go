package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

const headerKID = "kid"

type Client struct {
	activeKeyPair rsaKeyPair
	keyChain      rsaKeyChain
	expiration    time.Duration
}

func NewClient(cfg Config) (*Client, error) {
	var (
		activeKeyPair *rsaKeyPair
		keyChain      = make(rsaKeyChain, len(cfg.KeyPairs))
	)
	for _, kp := range cfg.KeyPairs {
		prv, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(kp.Private))
		if err != nil {
			return nil, fmt.Errorf("parsing RSA private key: %w", err)
		}

		pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(kp.Public))
		if err != nil {
			return nil, fmt.Errorf("parsing RSA public key: %w", err)
		}

		keyPair := rsaKeyPair{
			name:    kp.Name,
			private: prv,
			public:  pub,
		}

		if keyPair.name == cfg.ActiveKeyPair {
			activeKeyPair = &keyPair
		}

		keyChain[keyPair.name] = keyPair
	}

	if activeKeyPair == nil {
		return nil, fmt.Errorf("no such key pair found with name(%s)", cfg.ActiveKeyPair)
	}

	return &Client{
		activeKeyPair: *activeKeyPair,
		keyChain:      keyChain,
		expiration:    cfg.Expiration,
	}, nil
}

func (c *Client) ParseGoogleIDTokenUnverified(tokenString string) (*model.GoogleIDToken, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &googleIDTokenClaims{})
	if err != nil {
		return nil, fmt.Errorf("parsing google ID token: %w", err)
	}

	claims, ok := token.Claims.(*googleIDTokenClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected token claim type")
	}

	return claims.ToGoogleIDToken(), nil
}

func (c *Client) SignAppIDToken(appToken *model.AppIDToken) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, newAppClaims(appToken, c.expiration))
	token.Header["kid"] = c.activeKeyPair.name

	tokenString, err = token.SignedString(c.activeKeyPair.private)
	if err != nil {
		return "", fmt.Errorf("signing app ID token: %w", err)
	}

	return tokenString, nil
}

func (c *Client) VerifyAppIDTokenString(tokenString string) (*model.AppIDToken, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&appClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header[headerKID].(string)
			if !ok {
				return nil, fmt.Errorf("kid(%v) of token is not valid", token.Header[headerKID])
			}
			pair, ok := c.keyChain[kid]
			if !ok {
				return nil, fmt.Errorf("kid(%s) is not found from key chain", kid)
			}

			return pair.public, nil
		},
		jwt.WithIssuedAt(), jwt.WithIssuer(webMemoIssuer))
	switch {
	case err != nil:
		return nil, fmt.Errorf("parsing claim: %w", err)
	case !token.Valid:
		return nil, fmt.Errorf("invalid app ID token")
	}

	appClaims, ok := token.Claims.(*appClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected claim format")
	}

	appToken, err := appClaims.toAppIDToken()
	if err != nil {
		return nil, fmt.Errorf("converting to app ID token: %w", err)
	}

	return appToken, nil
}
