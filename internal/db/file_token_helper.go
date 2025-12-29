package db

func createToken() (token string, err error) {

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenBytes := make([]byte, 32)
	for i := range tokenBytes {
		tokenBytes[i] = chars[i%len(chars)]
	}
	return string(tokenBytes), nil
}
