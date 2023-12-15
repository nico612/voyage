package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestJWTAuth_GeneratorToken(t *testing.T) {

	key, _ := generateRandomKey(32)

	jwtAuth := NewJwtAuth([]byte(key))

	claims := NewClaims(&Options{
		Key:        key,
		Timeout:    1 * time.Hour,
		MaxRefresh: 2 * time.Hour,
	})
	claims.UserID = 1
	claims.Username = "zhangsan"

	tokenString, err := jwtAuth.GeneratorToken(claims)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(tokenString)

	// 解析
	cla, err := jwtAuth.ParseTokenString(tokenString)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("parse claimas = %v", cla)

}
