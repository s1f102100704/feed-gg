package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	dbgen "feed-gg/backend/internal/db"
	"feed-gg/backend/internal/httpapi"
	"feed-gg/backend/internal/riot"

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
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	riotClient := riot.NewClient(riotAPIKey)
	playerSearchHandler := httpapi.NewPlayerSearchHandler(riotClient)

	r := chi.NewRouter()
	r.Use(corsMiddleware)

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

	r.Post("/api/players/search", playerSearchHandler.Search)
	r.Options("/api/players/search", playerSearchHandler.Search)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
