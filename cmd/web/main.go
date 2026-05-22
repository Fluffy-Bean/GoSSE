package main

import (
	"context"
	"log/slog"
	"net/http"

	"git.leggy.dev/Fluffy/GoSSE/internal/handler"
	"git.leggy.dev/Fluffy/GoSSE/internal/routes"
	"git.leggy.dev/Fluffy/GoSSE/internal/sse"
)

const secret = "ManedWolves"

func main() {
	ctx := context.Background()

	s := sse.NewSSE(ctx)
	h := handler.NewHandler(s, secret)
	r := http.NewServeMux()

	routes.RegisterRoutes(h, r)

	slog.Info("Serving on :8080")
	http.ListenAndServe(":8080", r)
}
