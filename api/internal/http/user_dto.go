package http

import (
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type user struct {
	ID         uuid.UUID `json:"id"`
	UserType   string    `json:"userType"`
	Email      string    `json:"email"`
	UserName   string    `json:"userName"`
	GivenName  string    `json:"givenName,omitempty"`
	FamilyName string    `json:"familyName,omitempty"`
	PhotoURL   string    `json:"photoUrl,omitempty"`
	IssuedAt   time.Time `json:"issuedAt"`
	ExpireAt   time.Time `json:"expireAt"`
}

func (u *user) fromAppIDToken(token *model.AppIDToken) {
	u.ID = token.UserID
	u.UserType = string(token.UserType)
	u.Email = token.Email
	u.UserName = token.UserName
	u.GivenName = token.GivenName
	u.FamilyName = token.FamilyName
	u.PhotoURL = token.PhotoURL
	u.IssuedAt = token.IssuedAt
	u.ExpireAt = token.ExpireAt
}
