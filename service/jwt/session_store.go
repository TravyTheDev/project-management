package jwt

import (
	"database/sql"
	"fmt"
	"project-management/types"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{
		db: db,
	}
}

func (st *SessionStore) CreateSession(s *types.Session) (*types.Session, error) {
	stmt := `INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES (?, ?, ?, ?, ?)`

	_, err := st.db.Exec(stmt, s.ID, s.UserEmail, s.RefreshToken, s.IsRevoked, s.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	return s, nil
}

func (st *SessionStore) GetSession(id string) (*types.Session, error) {
	stmt := `SELECT * FROM sessions WHERE id = ?`

	rows, err := st.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	s := new(types.Session)
	for rows.Next() {
		err := rows.Scan(
			&s.ID,
			&s.UserEmail,
			&s.RefreshToken,
			&s.IsRevoked,
			&s.CreatedAt,
			&s.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (st *SessionStore) RevokeSession(id string) error {

	stmt := `UPDATE sessions SET is_revoked=1 WHERE id = ?`
	_, err := st.db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (st *SessionStore) DeleteSession(id string) error {
	_, err := st.db.Exec("DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}
