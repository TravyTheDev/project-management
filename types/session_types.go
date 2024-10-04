package types

import "time"

type SessionStore interface {
	CreateSession(*Session) (*Session, error)
	GetSession(string) (*Session, error)
	RevokeSession(string) error
	DeleteSession(string) error
}

type Session struct {
	ID           string    `json:"id"`
	UserEmail    string    `json:"user_email"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
