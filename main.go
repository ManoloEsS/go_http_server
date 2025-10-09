package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	//Root path to serve files from
	const filepathRoot = "."
	//default port for serving
	const port = "8080"

	//initialize config to share state
	cfg := &apiConfig{}

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
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	//handler for server health check
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	//handler for hit metrics check
	mux.HandleFunc("GET /admin/metrics", cfg.handlerRequestMetrics)
	//handler to reset metrics
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)
	//handler to validate chirp length
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	//listen and serve that blocks the log.Fatal server shutdown
	log.Fatal(server.ListenAndServe())
}

// api config to store a state
type apiConfig struct {
	fileServerHits atomic.Int32
}
