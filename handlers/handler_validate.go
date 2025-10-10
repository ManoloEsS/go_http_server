package handlers

import (
	"encoding/json"
	"go_http_server/internal/config"
	"net/http"
	"strings"
)

// function that validates a chirp and responds with a JSON or error
func HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	//chirp struct to use for decoding request body
	type chirp struct {
		Body string `json:"body"`
	}

	//response struct to respond without error
	type response struct {
		Body string `json:"cleaned_body"`
	}

	//decode response body and respond with error if not able to decode
	decoder := json.NewDecoder(r.Body)
	newChirp := chirp{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	//if chirp body is too long respond with error stating chirp is too long
	if len(newChirp.Body) > config.MaxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	cleanedChirp := profaneFilter(newChirp.Body)
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
