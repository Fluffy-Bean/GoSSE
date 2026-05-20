package handler

import (
	"git.leggy.dev/Fluffy/GoSSE/internal/sse"
)

type Handler struct {
	SSE *sse.SSE
}

func NewHandler(s *sse.SSE) *Handler {
	return &Handler{SSE: s}
}
