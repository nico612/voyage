package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomKey(size int) (string, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
