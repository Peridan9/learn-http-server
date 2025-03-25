package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/peridan9/learn-http-server/internal/auth"
)

func (cfg *APIConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type UserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		UserResponse
	}

	// decoding the request
	decoder := json.NewDecoder(r.Body)
	parameters := UserLogin{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	// checking there is an email and password
	if parameters.Email == "" || parameters.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", err)
		return
	}

	// getting the user
	user, err := cfg.DB.GetUserByEmail(r.Context(), parameters.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// checking the password
	err = auth.CheckPasswordHash(user.HashedPassword, parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// responding with the user
	respondWithJSON(w, http.StatusOK, response{UserResponse: UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}})
}
