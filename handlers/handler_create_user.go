package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	newResponseUser := responseUser{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	respondWithJSON(w, 201, newResponseUser)
}

// struct used to create the json response with appropriate fields from user created by cfg.Db.CreateUser
type responseUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
