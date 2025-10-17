package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ManoloEsS/go_http_server/server"
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
		server.RespondWithError(w, 500, "Couldn't decode chirp", err)
		return
	}

}
