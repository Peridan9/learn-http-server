package auth

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("no token provided")
	}
	splitAuth := strings.Split(tokenString, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	encodedKey := hex.EncodeToString(key)
	return encodedKey, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no API key provided")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
