package handlers

import (
	"go_http_server/internal/database"
	"sync/atomic"
)

// api config to store a state
type ApiConfig struct {
	fileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
}
