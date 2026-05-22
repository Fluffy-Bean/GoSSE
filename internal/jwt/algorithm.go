package jwt

import (
	"crypto"
	"crypto/hmac"
	"hash"
)

type Algorithm string

const (
	HS256 Algorithm = "HS256"
	HS384 Algorithm = "HS384"
	HS512 Algorithm = "HS512"
)

func sign(algorithm Algorithm, content, secret string) []byte {
	var hasher hash.Hash

	switch algorithm {
	case HS256:
		hasher = hmac.New(crypto.SHA256.New, []byte(secret))
	case HS384:
		hasher = hmac.New(crypto.SHA384.New, []byte(secret))
	case HS512:
		hasher = hmac.New(crypto.SHA512.New, []byte(secret))
	default:
		panic("unknown algorithm")
	}

	if _, err := hasher.Write([]byte(content)); err != nil {
		panic(err)
	}

	return hasher.Sum(nil)
}
