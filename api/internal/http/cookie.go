package http

import (
	"net/http"
	"time"
)

const cookieNameWebMemoToken = "wmToken"

func newWebMemoCookie(token string, ttl time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     cookieNameWebMemoToken,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(ttl),
	}
}
