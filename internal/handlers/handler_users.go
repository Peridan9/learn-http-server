package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/peridan9/learn-http-server/internal/auth"
	"github.com/peridan9/learn-http-server/internal/database"
)

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *APIConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type UserCreate struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		UserResponse
	}

	// decoding the request
	decoder := json.NewDecoder(r.Body)
	parameters := UserCreate{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	// checking there is an email and password
	if parameters.Email == "" || parameters.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	// hashing the password
	HashedPassword, err := auth.HashPassword(parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	// creating the user
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          parameters.Email,
		HashedPassword: HashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}

	// responding with the user
	respondWithJSON(w, http.StatusCreated, response{
		UserResponse: NewUserResponse(user),
	})
}
