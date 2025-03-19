package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	type Response struct {
		CleandBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	// checking if the chirp body length is longer then 140 characters
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respone := Response{
		CleandBody: cleanChirpFromBadWords(chirp.Body),
	}

	respondWithJSON(w, http.StatusOK, respone)
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
