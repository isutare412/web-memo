package handlers

import (
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// UserToWeb converts an AppIDToken to the generated User type.
func UserToWeb(token *model.AppIDToken) gen.User {
	return gen.User{
		ID:         token.UserID,
		UserType:   string(token.UserType),
		Email:      token.Email,
		UserName:   token.UserName,
		GivenName:  new(token.GivenName),
		FamilyName: new(token.FamilyName),
		PhotoURL:   new(token.PhotoURL),
		IssuedAt:   token.IssuedAt,
		ExpireAt:   token.ExpireAt,
	}
}
