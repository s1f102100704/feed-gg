package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	httpadapter "feed-gg/backend/internal/adapter/http"
	cacheinfra "feed-gg/backend/internal/infrastructure/cache"
	dbgen "feed-gg/backend/internal/infrastructure/db/sqlc"
	labelsinfra "feed-gg/backend/internal/infrastructure/labels"
	playersearchinfra "feed-gg/backend/internal/infrastructure/playersearch"
	"feed-gg/backend/internal/infrastructure/riot"
	"feed-gg/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverReadTimeout       = 15 * time.Second
	serverWriteTimeout      = 15 * time.Second
	serverIdleTimeout       = 60 * time.Second
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
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	queries := dbgen.New(sqlDB)
	labelsCache := cacheinfra.NewLabelsCache()
	labelsRepository := labelsinfra.NewRepository(queries)
	labelsUsecase := usecase.NewLabels(labelsCache, labelsRepository)
	labelsHandler := httpadapter.NewLabelsHandler(labelsUsecase)
	playerLabelsUsecase := usecase.NewPlayerLabels(labelsRepository)
	playerLabelsHandler := httpadapter.NewPlayerLabelsHandler(playerLabelsUsecase)
	playerSearchCache := cacheinfra.NewPlayerSearchCache()
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	if riotAPIKey == "" {
		log.Fatal("RIOT_API_KEY is not set")
	}
	riotClient := riot.NewClient(riotAPIKey)
	playerSearchRepository := playersearchinfra.NewRepository(sqlDB, queries, riotClient)
	playerSearchUsecase := usecase.NewPlayerSearch(
		playerSearchCache,
		playerSearchRepository,
		riotClient,
	)
	playerSearchHandler := httpadapter.NewPlayerSearchHandler(playerSearchUsecase)
	regionsHandler := httpadapter.NewRegionsHandler()

	r := chi.NewRouter()
	r.Use(corsMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			log.Printf("failed to write root response: %v", err)
		}
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := queries.Healthcheck(r.Context()); err != nil {
			http.Error(w, "db not ready", http.StatusServiceUnavailable)
			return
		}
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	})

	r.Get("/api/players/{region}/{gameName}/{tagLine}", playerSearchHandler.Search)
	r.Options("/api/players/{region}/{gameName}/{tagLine}", playerSearchHandler.Search)
	r.Get("/api/players/{puuid}/labels", playerLabelsHandler.List)
	r.Post("/api/players/{puuid}/labels", playerLabelsHandler.Vote)
	r.Options("/api/players/{puuid}/labels", playerLabelsHandler.List)
	r.Get("/api/labels", labelsHandler.List)
	r.Options("/api/labels", labelsHandler.List)
	r.Get("/api/regions", regionsHandler.List)
	r.Options("/api/regions", regionsHandler.List)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
	}

	log.Fatal(server.ListenAndServe())
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
