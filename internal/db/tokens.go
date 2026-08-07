package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

func VerifyAccessToken(token string) (bool, error) {
	var expiry string
	err := dbConn.QueryRow("SELECT expiry FROM accessTokens WHERE token = ?", token).Scan(&expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if expiry < time.Now().Format("2006-01-02 15:04:05") {
		return false, nil
	}

	return true, nil
}

func InvalidateAccessToken(token string) error {
	_, err := dbConn.Exec("DELETE FROM accessTokens WHERE token = ?", token)
	return err
}

func CreateAccessToken(tokenType string, expiry time.Time) (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	_, err = dbConn.Exec("INSERT INTO accessTokens (token, token_type, expiry) VALUES (?, ?, ?)", token, tokenType, expiry.Format("2006-01-02 15:04:05"))
	if err != nil {
		return "", err
	}

	return token, nil
}
