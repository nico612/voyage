package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

// Encrypt encrypts the plain text with bcrypt.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Sign issue a jwt token based on secretID, secretKey, iss and aud.
func Sign(secretID string, secretKey string, iss, aud string) string {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(0).Unix(),
		"aud": aud,
		"iss": iss,
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = secretID

	// Sign the token with the specified secret.
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}
