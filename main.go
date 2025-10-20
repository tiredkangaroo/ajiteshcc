package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/tiredkangaroo/ajiteshcc/bucket"
	"github.com/tiredkangaroo/ajiteshcc/env"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
	"github.com/tiredkangaroo/ajiteshcc/server"
)

func main() {
	if err := bucket.Init(); err != nil {
		slog.Error("bucket initialization", "error", err)
		return
	}
	conn, err := pgx.Connect(context.Background(), env.DefaultEnv.POSTGRES_CONNECTION_URI)
	if err != nil {
		slog.Error("database connection", "error", err)
		return
	}
	defer conn.Close(context.Background())
	slog.Info("database connected successfully")

	queries := db.New(conn)
	srv := &server.Server{Conn: conn, Queries: queries}
	if err := srv.Run(); err != nil {
		slog.Error("server run", "error", err)
		return
	}
}
