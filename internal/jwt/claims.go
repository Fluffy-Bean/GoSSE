package jwt

type Claims struct {
	Sub int64 `json:"sub"`
	Iat int64 `json:"iat"`
}
