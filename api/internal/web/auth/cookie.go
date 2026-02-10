package auth

import (
	"net/http"
	"time"
)

// CookieName is the name of the authentication token cookie.
const CookieName = "wmToken"

// NewWebMemoCookie creates a new HTTP cookie for the web-memo authentication token.
func NewWebMemoCookie(token string, ttl time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(ttl),
	}
}
