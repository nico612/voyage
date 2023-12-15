package auth

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {

	pwd, err := Encrypt("123456")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pwd)
}
