package handlers

import (
	"context"
	"encoding/json"
	"go_http_server/internal/config"
	"go_http_server/internal/database"
	"net/http"
	"strings"
)

// function that takes a request to /api/chirps and responds with a JSON or error response
func (cfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	//chirp struct to use for decoding request
	newRequestChirpParams := database.CreateChirpParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newRequestChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode chirp", err)
	}

	//validate chirp body length and filter profanity
	//if chirp body is too long respond with error stating chirp is too long
	if len(newRequestChirpParams.Body) > config.MaxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	newRequestChirpParams.Body = profaneFilter(newRequestChirpParams.Body)

	//Add chirp to database and return the struct
	validatedChirp, err := cfg.Db.CreateChirp(context.Background(), newRequestChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't save chirp to database", err)
	}

	//respond with success code and response instance
	respondWithJSON(w, 200, validatedChirp)

}

func profaneFilter(prefilter string) string {
	splitString := strings.Split(prefilter, " ")
	for i, word := range splitString {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			splitString[i] = "****"
		}
	}
	return strings.Join(splitString, " ")
}
