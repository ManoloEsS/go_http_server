package main

import (
	"database/sql"
	"go_http_server/handlers"
	"go_http_server/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//load .env file
	godotenv.Load()

	//get database connection url
	dbURL := os.Getenv("DB_URL")
	dbPlatform := os.Getenv("PLATFORM")

	//open
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect with database")
	}
	dbQueries := database.New(db)

	//Root path to serve files from
	const filepathRoot = "."
	//default port for serving
	const port = "8080"

	//initialize config to share state
	cfg := &handlers.ApiConfig{
		Db: dbQueries,
	}

	//initialize file server
	fileServer := http.FileServer(http.Dir(filepathRoot))

	//initialize multiplexer to handle requests
	mux := http.NewServeMux()

	//initialize http server struct
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	//main page handler, mapped to app and stripped to be "/"
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	//handler for server health check
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	//handler for hit metrics check
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerRequestMetrics)
	//handler to reset metrics
	mux.HandleFunc("POST /admin/reset", cfg.HandlerResetMetrics)
	//handler to validate chirp length
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirp)
	//handler to create a user
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	//listen and serve that blocks the log.Fatal server shutdown
	log.Fatal(server.ListenAndServe())
}
