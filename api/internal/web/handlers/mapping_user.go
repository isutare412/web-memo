package handlers

import (
	"github.com/samber/lo"

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
		GivenName:  lo.ToPtr(token.GivenName),
		FamilyName: lo.ToPtr(token.FamilyName),
		PhotoURL:   lo.ToPtr(token.PhotoURL),
		IssuedAt:   token.IssuedAt,
		ExpireAt:   token.ExpireAt,
	}
}
