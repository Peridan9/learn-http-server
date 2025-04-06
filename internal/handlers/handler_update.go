package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/peridan9/learn-http-server/internal/auth"
	"github.com/peridan9/learn-http-server/internal/database"
)

func (cfg *APIConfig) handlerUpdatePasswordAndEmail(w http.ResponseWriter, r *http.Request) {
	type UserUpdate struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		UserResponse
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	parameters := UserUpdate{}
	err = decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	if parameters.Email == "" && parameters.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email or password are required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	user, err := cfg.DB.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		ID:             userID,
		Email:          parameters.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		UserResponse: UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}
