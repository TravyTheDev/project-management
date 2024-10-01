package types

import "time"

type UserStore interface {
	CreateUser(User) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*UserRes, error)
	ChangePassword(string, string) error
}

type SessionStore interface {
	CreateSession(*Session) (*Session, error)
	GetSession(string) (*Session, error)
	RevokeSession(string) error
	DeleteSession(string) error
}

type User struct {
	ID        int64      `db:"id"`
	Username  string     `db:"username"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type RegisterPayload struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=3,max=130"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type UserRes struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Session struct {
	ID           string    `db:"id"`
	UserEmail    string    `db:"user_email"`
	RefreshToken string    `db:"refresh_token"`
	IsRevoked    bool      `db:"is_revoked"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}

type CreateRoomReq struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type NumbersConfirmReq struct {
	Numbers int `json:"numbers"`
}

type PasswordResetReq struct {
	Password        string `json:"password" validate:"required,min=3,max=130"`
	PasswordConfirm string `json:"passwordConfirm"`
}
