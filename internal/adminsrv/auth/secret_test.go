package auth

import (
	"fmt"
	"testing"
)

func TestGenerateRandomKey(t *testing.T) {

	// 生成 256位（32字节）的密钥
	ker, err := generateRandomKey(32)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ker)

}
