package handlers

import (
	"fmt"
	"net/http"

	"github.com/ManoloEsS/go_http_server/server"
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
func (cfg *ApiConfig) HandlerResetUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		server.RespondWithError(w, http.StatusForbidden, "Reset is only allowed in env environment", nil)
		return
	}
	cfg.fileServerHits.Store(0)
	err := cfg.Db.DeleteAllUsers(r.Context())
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Failed to reset database", err)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users have been deleted and hits reset to 0"))
}
