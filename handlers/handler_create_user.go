package handlers

import (
	"encoding/json"
	"net/http"
)

// Method of ApiConfig. takes a http.ResponseWriter and a http.Request and creates
// a new user in the database and responds with the user data in JSON
func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userEmail := email{}
	err := decoder.Decode(&userEmail)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode user email", err)
	}

	newUser, err := cfg.Db.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		respondWithError(w, 500, "Couldn't create user", err)
	}

	respondWithJSON(w, 201, newUser)
}
