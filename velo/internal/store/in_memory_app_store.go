package store

import (
	"fmt"
	"sync"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
)

type InMemoryAppStore struct {
	mu    sync.RWMutex
	Items []models.Item
}

func NewInMemoryAppStore() *InMemoryAppStore {
	return &InMemoryAppStore{
		sync.RWMutex{},
		[]models.Item{},
	}
}

func (s *InMemoryAppStore) GetItem(id string) (models.Item, error) {
	return models.Item{}, fmt.Errorf("item not found: %s", id)
}

func (s *InMemoryAppStore) GetItems() ([]models.Item, error) {
	return []models.Item{}, nil
}

func (s *InMemoryAppStore) DeleteItem(id string) error {
	return fmt.Errorf("item not found when trying to delete it: %s", id)
}

func (s *InMemoryAppStore) UpdateItem(id string, updates map[string]any) (models.Item, error) {
	return models.Item{}, nil
}
