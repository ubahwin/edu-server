package model

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID                    uuid.UUID
	Scope                 SessionScope
	AccessTokenOfVdovinID string
	AccessToken           string
	RefreshToken          string
	AccessTokenTTL        time.Duration
	UpdatedAt             time.Time
	CreatedAt             time.Time
}
