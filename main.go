package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/tiredkangaroo/ajiteshcc/env"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
	"github.com/tiredkangaroo/ajiteshcc/server"
)

func main() {
	conn, err := pgx.Connect(context.Background(), env.DefaultEnv.POSTGRES_CONNECTION_URI)
	if err != nil {
		slog.Error("database connection", "error", err)
		return
	}
	defer conn.Close(context.Background())
	slog.Info("database connected successfully")

	queries := db.New(conn)
	srv := &server.Server{Queries: queries}
	if err := srv.Run(); err != nil {
		slog.Error("server run", "error", err)
		return
	}
}
