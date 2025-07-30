package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
)

// Handler manages the API endpoints
type Handler struct {
	store models.AppStore
}

func NewHandler(store models.AppStore) *Handler {
	return &Handler{store: store}
}

// ServeHTTP handles all API requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Route to the appropriate handler based on the path
	path := r.URL.Path

	if path == "/items" {
		h.handleItems(w, r)
		return
	}

	if strings.HasPrefix(path, "/items/") {
		h.handleItem(w, r)
		return
	}

	// Handle unknown paths
	w.WriteHeader(http.StatusNotFound)
}

// handleItem processes requests for individual items
func (h *Handler) handleItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")

	switch r.Method {
	case http.MethodGet:
		h.getItem(w, id)
	case http.MethodPatch:
		h.updateItem(w, r, id)
	case http.MethodDelete:
		h.deleteItem(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getItem retrieves a single item
func (h *Handler) getItem(w http.ResponseWriter, id string) {
	item, err := h.store.GetItem(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// updateItem updates an existing item
func (h *Handler) updateItem(w http.ResponseWriter, r *http.Request, id string) {
	var updates map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedItem, err := h.store.UpdateItem(id, updates)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedItem)
}

// deleteItem removes an item
func (h *Handler) deleteItem(w http.ResponseWriter, id string) {
	err := h.store.DeleteItem(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleItems processes requests for collections of items
func (h *Handler) handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getItems(w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getItems retrieves all items
func (h *Handler) getItems(w http.ResponseWriter) {
	items, err := h.store.GetItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

// respondWithJSON sends a JSON response with the given status code
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
