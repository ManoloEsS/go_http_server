package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go_http_server/internal/config"
	"go_http_server/internal/database"
	"net/http"
	"server"
	"strings"
	"time"

	"github.com/google/uuid"
)

// function that takes a request to /api/chirps, creates a new Chirp saved at table Chirps and responds with a JSON struct of the Chirp data or error response
func (cfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	//chirp struct to use for decoding request
	type ChirpParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	newRequestChirpParams := ChirpParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newRequestChirpParams)
	if err != nil {
		server.RespondWithJSON(w, http.StatusInternalServerError, "Couldn't decode chirp", err)
		return
	}

	//validate chirp body length and filter profanity
	filteredChirpBody, err := validateChirp(newRequestChirpParams.Body)
	if err != nil {
		server.RespondWithJSON(w, http.StatusBadRequest, err.Error(), err)
	}

	validatedChirpData := database.CreateChirpParams{
		Body:   filteredChirpBody,
		UserID: newRequestChirpParams.UserID,
	}
	//Add chirp to database and return the struct
	validatedChirp, err := cfg.Db.CreateChirp(context.Background(), validatedChirpData)
	if err != nil {
		server.RespondWithJSON(w, http.StatusInternalServerError, "couldn't save chirp to database", err)
		return
	}

	resp := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}{
		ID:        validatedChirp.ID,
		CreatedAt: validatedChirp.CreatedAt,
		UpdatedAt: validatedChirp.CreatedAt,
		Body:      validatedChirp.Body,
		UserID:    validatedChirp.UserID,
	}

	//respond with success code and response instance
	server.RespondWithJSON(w, 201, resp)
}

func validateChirp(body string) (string, error) {
	if len(body) > config.MaxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	filteredBody := profaneFilter(body, badWords)

	return filteredBody, nil

}
func profaneFilter(prefilter string, profanity map[string]struct{}) string {
	splitString := strings.Split(prefilter, " ")
	for i, word := range splitString {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			splitString[i] = "****"
		}
	}
	return strings.Join(splitString, " ")
}

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}
