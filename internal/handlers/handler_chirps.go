package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peridan9/learn-http-server/internal/auth"
	"github.com/peridan9/learn-http-server/internal/database"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *APIConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get chirp", err)
		return
	}

	response := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *APIConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}

	response := make([]ChirpResponse, 0, len(chirps))
	for _, chirp := range chirps {
		response = append(response, ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}
func (cfg *APIConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type Chirp struct {
		Body string `json:"body"`
	}

	// TODO: implement a middleware to check if the user is authenticated
	// check if the user is authenticated
	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// validate the token
	UserID, err := auth.ValidateJWT(userToken, cfg.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	// checking if user exists
	_, err = cfg.DB.GetUserByID(r.Context(), UserID)
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
		UserID: UserID,
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
