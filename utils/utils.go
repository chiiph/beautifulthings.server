package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func S25664(s string) string {
	return base64.URLEncoding.EncodeToString(S256(s))
}

func S256(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}
