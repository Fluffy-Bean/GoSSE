package handler

import (
	"fmt"
	"net/http"
	"time"

	"git.leggy.dev/Fluffy/GoSSE/internal/jwt"
	"git.leggy.dev/Fluffy/GoSSE/internal/sse"
)

// On HTTPS connections you want to prefix your cookies with __Secure
const (
	cookieName     = "I_Heart_Maned_Wolves"
	cookieLifespan = 60 * time.Minute
)

type Handler struct {
	SSE       *sse.SSE
	secret    string
	startedAt time.Time
}

func NewHandler(s *sse.SSE, secret string) *Handler {
	return &Handler{
		SSE:       s,
		secret:    secret,
		startedAt: time.Now().UTC(),
	}
}

func (h *Handler) SetToken(w http.ResponseWriter, id int64) error {
	now := time.Now().UTC()

	claims := jwt.Claims{
		Sub: id,
		Iat: now.Unix(),
	}

	token, err := jwt.New(jwt.HS256, h.secret).Encode(claims)
	if err != nil {
		return fmt.Errorf("encode token: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		HttpOnly: true,
		Expires:  now.Add(cookieLifespan),
	})

	return nil
}

func (h *Handler) GetToken(r *http.Request) (*jwt.Claims, error) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == cookieName {
			cookie = c

			break
		}
	}

	if cookie == nil {
		return nil, fmt.Errorf("no cookie set")
	}

	claims, err := jwt.New(jwt.HS256, h.secret).VerifyAndDecode(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("decoding jwt token: %w", err)
	}

	if claims.Iat < h.startedAt.UTC().Unix() {
		return nil, fmt.Errorf("token is too old")
	}

	return claims, nil
}
