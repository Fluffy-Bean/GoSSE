package routes

import (
	"encoding/base64"
	"encoding/json"
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

		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("cache-control", "no-cache")
		w.Header().Set("connection", "keep-alive")

		w.WriteHeader(http.StatusOK)

		if err := rc.Flush(); err != nil {
			slog.Error("flush headers", "error", err)

			return
		}

		h.SSE.Subscribe(conn)
		defer h.SSE.Unsubscribe(conn)

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

				data := base64.StdEncoding.EncodeToString([]byte(message))

				w.Write([]byte("event: message\n"))
				w.Write([]byte("data: " + data + "\n"))
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
		var input struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			slog.Error("decode chat message", "error", err)

			return
		}

		go h.SSE.Broadcast(input.Message)

		w.WriteHeader(http.StatusAccepted)
	}
}
