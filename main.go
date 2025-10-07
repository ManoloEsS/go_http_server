package main

import (
	"fmt"
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
	mux.HandleFunc("/healthz", handlerReadiness)
	//handler for hit metrics check
	mux.HandleFunc("/metrics", cfg.handlerRequestMetrics)
	//handler to reset metrics
	mux.HandleFunc("/reset", cfg.handlerResetMetrics)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	//listen and serve that blocks the log.Fatal server shutdown
	log.Fatal(server.ListenAndServe())
}

// function that writes the response for health check
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// function that writes the response for metrics check
func (cfg *apiConfig) handlerRequestMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())))
}

// function that writes the response for reset metrics
func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics have been reset"))
}

// middleware for handlers that add a hit to the fileserver metric
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// api config to store a state
type apiConfig struct {
	fileServerHits atomic.Int32
}
