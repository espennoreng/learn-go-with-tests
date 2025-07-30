package velo

import (
	"net/http"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
	"github.com/espennoreng/learn-go-with-tests/velo/pkg/api"
)

// AppServer is the main server that handles HTTP requests
type AppServer struct {
	store models.AppStore
	http.Handler
}

// NewAppServer creates and configures a new server instance
func NewAppServer(store models.AppStore) *AppServer {
	s := &AppServer{
		store: store,
	}

	// Create API handlers with the provided store
	apiHandler := api.NewHandler(store)

	// Configure routing
	router := http.NewServeMux()
	router.Handle("/items/", apiHandler)
	router.Handle("/items", apiHandler)

	s.Handler = router
	return s
}
