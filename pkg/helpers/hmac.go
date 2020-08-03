package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func HMACSHA256(data []byte, key string) string {
	hm := hmac.New(sha256.New, []byte(key))

	hm.Write(data)

	return fmt.Sprintf("%x", hm.Sum(nil))
}
