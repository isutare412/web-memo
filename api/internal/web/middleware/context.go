package middleware

import (
	"context"
	"net/http"

	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/log"
)

type contextBagKey struct{}

// ContextBag holds request-scoped data that is passed through middleware chain.
type ContextBag struct {
	Passport *Passport
}

// Passport holds the authentication token and its parsed form.
type Passport struct {
	TokenString string
	Token       *model.AppIDToken
}

// WithContextBag initializes a ContextBag in the request context.
func WithContextBag(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		bag := &ContextBag{}
		ctx := context.WithValue(r.Context(), contextBagKey{}, bag)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetContextBag retrieves the ContextBag from the request context.
func GetContextBag(ctx context.Context) (*ContextBag, bool) {
	bag, ok := ctx.Value(contextBagKey{}).(*ContextBag)
	return bag, ok
}

// ExtractPassport retrieves the Passport from the request context's ContextBag.
func ExtractPassport(ctx context.Context) (*Passport, bool) {
	bag, ok := GetContextBag(ctx)
	if !ok || bag.Passport == nil {
		return nil, false
	}
	return bag.Passport, true
}

// WithLogAttrContext injects a log attribute context into the request context.
func WithLogAttrContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = log.WithAttrContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
