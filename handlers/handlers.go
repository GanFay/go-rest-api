package handlers

import "github.com/jackc/pgx/v5"

type Handler struct {
	DB *pgx.Conn
}

func New(db *pgx.Conn) *Handler {
	return &Handler{DB: db}
}
