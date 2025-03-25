package handlers

import (
	"sync/atomic"

	"github.com/peridan9/learn-http-server/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	SecretKey      string
}
