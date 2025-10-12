package handlers

import (
	"encoding/json"
	"go_http_server/internal/config"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// function that takes a request to /api/chirps and responds with a JSON or error response
func HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	//chirp struct to use for decoding request
	type chirp struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	requestData := chirp{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode chirp", err)
	}

	//if chirp body is too long respond with error stating chirp is too long
	if len(requestData.Body) > config.MaxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	cleanedChirp := profaneFilter(requestData.Body)
	//respond with success code and response instance
	respondWithJSON(w, 200, response{
		Body: cleanedChirp,
	})

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
