package handlers

import (
	"net/http"

	"github.com/ManoloEsS/go_http_server/server"
)

func (cfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Db.GetAllChirps(r.Context())
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, chirps)
}
