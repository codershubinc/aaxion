package db

import "crypto/rand"

func createToken() (token string, err error) {

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_~"
	tokenBytes := make([]byte, 32)
	randNum := make([]byte, 32)
	_, err = rand.Read(randNum)
	if err != nil {
		return "", err
	}
	for i := range 32 {
		tokenBytes[i] = chars[int(randNum[i])%len(chars)]
	}
	return string(tokenBytes), nil
}
