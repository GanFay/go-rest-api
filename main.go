package main

import (
	"api/handlers"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}

	url := os.Getenv("DB_URL")

	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *pgx.Conn, ctx context.Context) {
		err = db.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(db, context.Background())
	log.Println("Connected to database")

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			name       TEXT        NOT NULL,
			password   TEXT        NOT NULL,
			gmail      TEXT        NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatal("create table error:", err)
	}
	log.Println("TABLE CHECK/CREATED OK")
	var dbName string
	if err = db.QueryRow(context.Background(), "select current_database()").Scan(&dbName); err != nil {
		log.Fatal("cannot get db name:", err)
	}
	log.Println("CONNECTED TO DB:", dbName)

	h := handlers.New(db)

	//роутер
	http.HandleFunc("/health", handlers.Health)
	http.HandleFunc("/ping", handlers.Ping)
	http.HandleFunc("/users", h.Users)
	http.HandleFunc("/users/delete", h.Delete)

	addr := ":8080"
	log.Printf("Listening on %s\n", addr)
	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
