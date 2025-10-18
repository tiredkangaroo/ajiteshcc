package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/tiredkangaroo/ajiteshcc/env"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
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
	if err := server(queries); err != nil {
		slog.Error("server", "error", err)
		return
	}
}
