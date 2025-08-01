package api

import (
	"encoding/json"
	"fmt"
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

	if strings.HasPrefix(path, "/sessions/"){
		h.handleSessions(w, r)
		return
	}

	if strings.HasPrefix(path, "/users/"){
		h.handleUser(w, r)
		return
	}

	if strings.HasPrefix(path, "/users"){ 
		h.handleUsers(w, r)
		return
	}

	// Handle unknown paths
	w.WriteHeader(http.StatusNotFound)
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodPost:
		 h.createUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req CreateUserRequest
	// bad json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate the data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// valid data
	createdUser, err := h.store.CreateUser(models.CreateUserInput{
		Name: req.Name,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/users/%s", createdUser.ID))
	respondWithJSON(w, http.StatusCreated, createdUser)

}

func (h *Handler) handleUser(w http.ResponseWriter, r *http.Request){
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	switch r.Method{
	case http.MethodGet:
		h.getUser(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getUser(w http.ResponseWriter, id string){
	user, err := h.store.GetUser(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
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
	case http.MethodPost:
		h.createItem(w, r)
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

// createItem creates a item
func (h *Handler) createItem(w http.ResponseWriter, r *http.Request){
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req CreateItemRequest
	// bad json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate the data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// valid data
	createdItem, err := h.store.CreateItem(models.CreateItemInput{
		Name: req.Name,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/items/%s", createdItem.ID))
	respondWithJSON(w, http.StatusCreated, createdItem)

}


// handleSessions processes requests for sessions
func (h *Handler) handleSessions(w http.ResponseWriter, r *http.Request){
	id := strings.TrimPrefix(r.URL.Path, "/sessions/")

	switch r.Method{
	case http.MethodGet:
		h.getSession(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getSession(w http.ResponseWriter, id string){
	
	session, err := h.store.GetSession(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, session)
}



// respondWithJSON sends a JSON response with the given status code
func respondWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
