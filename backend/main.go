package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	dbgen "feed-gg/backend/internal/db"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	queries := dbgen.New(sqlDB)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := queries.Healthcheck(r.Context()); err != nil {
			http.Error(w, "db not ready", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ok"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
