package users

import (
	"database/sql"
	"fmt"
	"project-management/types"
	"strings"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) CreateUser(user types.User) (*types.User, error) {
	stmt := `INSERT INTO users (username, email, password, is_admin) VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(stmt, user.Username, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	u, err := s.GetUserByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting created user")
	}

	return u, nil
}

func (s *UserStore) GetUserByEmail(email string) (*types.User, error) {
	stmt := `SELECT * FROM users WHERE email = ?`

	rows, err := s.db.Query(stmt, email)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserStore) GetUserByID(id int) (*types.UserRes, error) {
	stmt := `SELECT * FROM users WHERE id = ?`

	rows, err := s.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	u := &types.UserRes{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return u, nil
}

func (s *UserStore) ChangePassword(email string, password string) error {
	stmt := `UPDATE users SET password = ? WHERE email = ?`
	_, err := s.db.Exec(stmt, password, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) SearchUser(text string) ([]*types.UserRes, error) {
	users := make([]*types.UserRes, 0)
	splitStr := strings.Split(text, " ")
	for _, search := range splitStr {
		if len(search) != 0 {
			stmt := `SELECT * FROM users WHERE username LIKE ?`
			rows, err := s.db.Query(stmt, `%`+search+`%`)
			if err != nil {
				return nil, err
			}
			for rows.Next() {
				u, err := scanRowsIntoUser(rows)
				if err != nil {
					return nil, err
				}
				user := &types.UserRes{
					ID:       u.ID,
					Username: u.Username,
					Email:    u.Email,
				}
				users = append(users, user)
			}
		}
	}
	return users, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
