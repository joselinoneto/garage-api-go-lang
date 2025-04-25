package models

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Create(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	_, err = m.DB.Exec(stmt, username, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) Authenticate(username, password string) (*User, error) {
	var user User

	stmt := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := m.DB.QueryRow(stmt, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
} 