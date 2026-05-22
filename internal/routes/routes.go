package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"git.leggy.dev/Fluffy/GoSSE/internal/handler"
	"git.leggy.dev/Fluffy/GoSSE/internal/sse"
)

func RegisterRoutes(h *handler.Handler, r *http.ServeMux) {
	r.HandleFunc("GET /", homeGet(h))
	r.HandleFunc("GET /connect", connGet(h))
	r.HandleFunc("POST /message", messagePost(h))
}

func homeGet(h *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			slog.Error("parse template", "error", err)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, nil)
		if err != nil {
			slog.Error("execute template", "error", err)

			return
		}
	}
}

func connGet(h *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rc := http.NewResponseController(w)
		conn := sse.NewConnection()

		// Restore user ID on reconnect
		claims, _ := h.GetToken(r)
		if claims != nil {
			conn.ID = claims.Sub
		}

		h.SSE.Subscribe(conn)
		defer h.SSE.Unsubscribe(conn)

		// Assign new JWT token if not already had one
		if claims == nil {
			if err := h.SetToken(w, conn.ID); err != nil {
				slog.Error("set token", "error", err)

				return
			}
		}

		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("cache-control", "no-cache")
		w.Header().Set("connection", "keep-alive")

		w.WriteHeader(http.StatusOK)

		if err := rc.Flush(); err != nil {
			slog.Error("flush headers", "error", err)

			return
		}

		rc.SetWriteDeadline(time.Now().Add(h.SSE.Heartbeat * 2))

		for {
			select {
			case _, ok := <-conn.Heartbeat:
				if !ok {
					return
				}

				rc.SetWriteDeadline(time.Now().Add(h.SSE.Heartbeat * 2))
				w.Write([]byte(":\n\n"))

				if err := rc.Flush(); err != nil {
					slog.Error("flush heartbeat", "error", err)

					return
				}

			case message, ok := <-conn.Messages:
				if !ok {
					return
				}

				userBase64 := base64Encode(fmt.Sprintf("%d", message.UserID))
				messageBase64 := base64Encode(message.Message)

				w.Write([]byte("event: message\n"))
				w.Write([]byte("data: " + userBase64 + "." + messageBase64 + "\n"))
				w.Write([]byte("\n\n"))

				if err := rc.Flush(); err != nil {
					slog.Error("flush event", "error", err)

					return
				}

			case <-ctx.Done():
				return
			}
		}
	}
}

func messagePost(h *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := h.GetToken(r)
		if err != nil {
			slog.Error("get token", "error", err)

			return
		}

		var input struct {
			UserID  int64
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			slog.Error("decode chat message", "error", err)

			return
		}

		input.UserID = claims.Sub

		go h.SSE.Broadcast(sse.Message(input))

		w.WriteHeader(http.StatusAccepted)
	}
}

func base64Encode(src string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(src))
}
