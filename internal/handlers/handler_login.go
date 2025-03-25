package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/peridan9/learn-http-server/internal/auth"
)

func (cfg *APIConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type UserLogin struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiersInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		UserResponse
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	// creating the token
	token, err := auth.MakeJWT(user.ID, cfg.SecretKey, clampTokenExpiry(parameters.ExpiersInSeconds))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}

	// responding with the user
	respondWithJSON(w, http.StatusOK, response{UserResponse: UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	},
		Token: token,
	})
}

// clampTokenExpiry ensures the token expiry is between 1 and 3600 seconds
func clampTokenExpiry(seconds int) time.Duration {
	if seconds <= 0 || seconds > 3600 {
		seconds = 3600
	}
	return time.Duration(seconds) * time.Second
}
