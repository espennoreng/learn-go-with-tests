package velo

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
)

type AppStore interface {
	GetItem(id string) (models.Item, error)
	GetItems() ([]models.Item, error)
	DeleteItem(id string) error
	UpdateItem(id string, update map[string]any) (models.Item, error)
}

type AppServer struct {
	store AppStore
	http.Handler
}

func NewAppServer(store AppStore) *AppServer {
	s := new(AppServer)
	s.store = store
	router := http.NewServeMux()

	router.Handle("/items/", http.HandlerFunc(s.itemHandler))
	router.Handle("/items", http.HandlerFunc(s.itemsHandler))

	s.Handler = router
	return s
}

func (a *AppServer) itemHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")

	switch r.Method {
	case http.MethodGet:
		item, err := a.store.GetItem(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)

	case http.MethodPatch:
		var updates map[string]any
		err := json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedItem, err := a.store.UpdateItem(id, updates)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedItem)

	case http.MethodDelete:
		err := a.store.DeleteItem(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}

}

func (a *AppServer) itemsHandler(w http.ResponseWriter, r *http.Request) {

	items, err := a.store.GetItems()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
