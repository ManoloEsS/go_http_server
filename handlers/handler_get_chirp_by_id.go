package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.Db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
