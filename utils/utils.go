package utils

import "crypto/sha256"

func HashPassword(password string) string {

	h := sha256.New()
	h.Write([]byte(password))
	bs := h.Sum(nil)

	return string(bs)
}
