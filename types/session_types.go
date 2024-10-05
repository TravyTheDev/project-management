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
	UserEmail    string    `json:"userEmail"`
	RefreshToken string    `json:"refreshToken"`
	IsRevoked    bool      `json:"isRevoked"`
	CreatedAt    time.Time `json:"createdAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
