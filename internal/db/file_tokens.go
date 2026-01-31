package db

import "log"

func CreateFileShareTempToken(filePath string) (token string, err error) {
	tk, err := createToken()
	if err != nil {
		log.Println("err", err)
		return "", err
	}
	_, err = GetDB().Exec("	INSERT INTO tokens (token, token_type, file_path, expiry) VALUES (?, ?, ?, datetime('now', '+1 hour'))",
		tk, "file_share", filePath)
	if err != nil {
		log.Println("err  for db", err)
		return "", err
	}
	return tk, nil
}

func ValidateFileShareToken(token string) (filePath string, err error) {
	row := GetDB().QueryRow("SELECT file_path FROM tokens WHERE token = ? AND token_type = ? AND expiry > datetime('now')",
		token, "file_share")
	err = row.Scan(&filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func RevokeFileShareToken(token string) error {
	_, err := GetDB().Exec("DELETE FROM tokens WHERE token = ? AND token_type = ?",
		token, "file_share")
	return err
}
