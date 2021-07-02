package weapp

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := "12345678901234567890123456789012"
	str, err := Encrypt(key, "50980")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)

	dStr, err := Decrypt(key, str)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dStr)
}
