package jwt

import (
	"encoding/base64"
)

func base64Encode(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

func base64Decode(src string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(src)
}
