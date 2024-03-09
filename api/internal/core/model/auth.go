package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
)

type GoogleTokenResponse struct {
	AccessToken    string
	AccessTokenTTL time.Duration
	IDToken        string
	Scope          string
	TokenType      string
}

type GoogleIDToken struct {
	IssuedAt   time.Time
	ExpiresAt  time.Time
	Subject    string
	Email      string
	Name       string
	FamilyName string
	GivenName  string
	PictureURL string
}

type AppIDToken struct {
	UserID     uuid.UUID
	UserType   enum.UserType
	Email      string
	UserName   string
	FamilyName string
	GivenName  string
	PhotoURL   string
}

func (t *AppIDToken) CanWriteMemo(memo *ent.Memo) bool {
	if t == nil {
		return false
	}

	if t.UserType == enum.UserTypeOperator {
		return true
	}

	if memo.OwnerID == t.UserID {
		return true
	}

	if _, ok := lo.Find(
		memo.Edges.Collaborations,
		func(c *ent.Collaboration) bool { return c.UserID == t.UserID && c.Approved }); ok {
		return true
	}

	return false
}

func (t *AppIDToken) CanReadMemo(memo *ent.Memo) bool {
	if t.CanWriteMemo(memo) {
		return true
	}
	return memo.IsPublished
}

func (t *AppIDToken) IsOwner(memo *ent.Memo) bool {
	if t == nil {
		return false
	}

	if t.UserType == enum.UserTypeOperator {
		return true
	}

	return memo.OwnerID == t.UserID
}
