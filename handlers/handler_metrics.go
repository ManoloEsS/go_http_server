package handlers

import (
	"fmt"
	"net/http"
)

// function that writes the response for metrics check
func (cfg *ApiConfig) HandlerRequestMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
	  <body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`, cfg.fileServerHits.Load())))
}

// function that writes the response for reset metrics
func (cfg *ApiConfig) HandlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics have been reset"))
}
