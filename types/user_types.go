package types

import "time"

type UserStore interface {
	CreateUser(User) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*UserRes, error)
	ChangePassword(string, string) error
	SearchUser(string) ([]*UserRes, error)
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RegisterPayload struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=3,max=130"`
	PasswordConfirm string `json:"passwordConfirm"`
	IsAdmin         bool   `json:"isAdmin"`
}

type UserRes struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SearchReq struct {
	Text string `json:"text"`
}
