package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

func CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = dbConn.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hashedPassword))
	return err
}

func AuthenticateUser(username, password string) (string, error) {
	var user User
	err := dbConn.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")

		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Create token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	_, err = dbConn.Exec("INSERT INTO auth_tokens (user_id, token) VALUES (?, ?)", user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token string) (bool, error) {
	var id int
	err := dbConn.QueryRow("SELECT id FROM auth_tokens WHERE token = ?", token).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func InvalidateToken(token string) error {
	_, err := dbConn.Exec("DELETE FROM auth_tokens WHERE token = ?", token)
	return err
}

func HasUsers() (bool, error) {
	var count int
	err := dbConn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
