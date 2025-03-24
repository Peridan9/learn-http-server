package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peridan9/learn-http-server/internal/database"
)

func (cfg *APIConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type Chirp struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	// checking if user exists
	_, err = cfg.DB.GetUserByID(r.Context(), uuid.MustParse(chirp.UserID))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User does not exist", err)
		return
	}

	// checking if the chirp body length is longer then 140 characters
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp.Body = cleanChirpFromBadWords(chirp.Body)

	newChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: uuid.MustParse(chirp.UserID),
		Body:   chirp.Body,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respone := ChirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, respone)

}

func cleanChirpFromBadWords(chirp string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(chirp, " ") // Split by spaces

	for i, word := range words {
		// Normalize the word to lowercase for comparison
		cleanedWord := strings.ToLower(word)

		// Check if the word (case-insensitive) is in the bad words list
		for _, badWord := range badWords {
			if cleanedWord == badWord {
				words[i] = "****" // Replace the bad word
				break
			}
		}
	}

	return strings.Join(words, " ") // Join words back into a sentence
}
