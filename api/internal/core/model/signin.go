package model

// GoogleSignInResult is the outcome of a completed Google sign-in.
type GoogleSignInResult struct {
	// RedirectURL is where the user agent must be redirected next.
	RedirectURL string
	AppToken    string
	// SetCookie is true if AppToken must be delivered as the web cookie. It is
	// false for CLI logins, which receive the token in RedirectURL instead.
	SetCookie bool
}
