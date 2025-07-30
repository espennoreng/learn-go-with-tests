package testutils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
)

type StubAppStore struct {
	mu    sync.RWMutex
	Items []models.Item
}

func NewStubAppStore() *StubAppStore{
	return &StubAppStore{}
}

func (s *StubAppStore) GetItem(id string) (models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.Items {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Item{}, fmt.Errorf("item not found: %s", id)
}

func (s *StubAppStore) GetItems() ([]models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Items, nil
}

func (s *StubAppStore) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, item := range s.Items {
		if item.ID == id {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item not found when trying to delete %s from %v", id, s.Items)
}

func (s *StubAppStore) UpdateItem(id string, updates map[string]any) (models.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Find the item
	for i, item := range s.Items {
		if item.ID == id {
			// Apply updates
			if name, ok := updates["Name"].(string); ok {
				s.Items[i].Name = name
			}
			if externalID, ok := updates["ExternalID"].(string); ok {
				s.Items[i].ExternalID = externalID
			}
			if orgID, ok := updates["OrgID"].(string); ok {
				s.Items[i].OrgID = orgID
			}
			if isActive, ok := updates["IsActive"].(string); ok {
				s.Items[i].IsActive = isActive
			}
			if createdBy, ok := updates["CreatedBy"].(string); ok {
				s.Items[i].CreatedBy = createdBy
			}
			if deletedAt, ok := updates["DeletedAt"].(string); ok {
				s.Items[i].DeletedAt = deletedAt
			}

			return s.Items[i], nil
		}
	}

	return models.Item{}, fmt.Errorf("item not found: %s", id)
}

type ErrorStore struct {
    models.AppStore
    ShouldErrorOnGetItems bool
}

func (s *ErrorStore) GetItems() ([]models.Item, error) {
    if s.ShouldErrorOnGetItems {
        return nil, fmt.Errorf("forced error for testing")
    }
    return s.AppStore.GetItems()
}



func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("response status is wrong, got %d, want %d", got, want)
	}
}

func AssertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func AssertValidJSON(t testing.TB, body *bytes.Buffer, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unable to parse response from server %q, '%v'", body, err)
	}
}


func MakeRequest(t testing.TB, server http.Handler, method, url string, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	return res
}

func AssertItemsContain(t testing.TB, items []models.Item, expectedIDs ...string) {
	t.Helper()

	if len(items) != len(expectedIDs) {
		t.Errorf("expected %d items, got %d", len(expectedIDs), len(items))
	}

	itemMap := make(map[string]models.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	for _, id := range expectedIDs {
		if _, ok := itemMap[id]; !ok {
			t.Errorf("expected item with ID %s not found", id)
		}
	}
}

func CreateTestItems() []models.Item {
	now := time.Now()
	return []models.Item{
		{
			ID:         "item-001",
			Name:       "First Test Item",
			ExternalID: "ext-001",
			OrgID:      "org-123",
			IsActive:   "true",
			CreatedAt:  now,
			CreatedBy:  "test-user",
			DeletedAt:  "",
		},
		{
			ID:         "item-002",
			Name:       "Second Test Item",
			ExternalID: "ext-002",
			OrgID:      "org-123",
			IsActive:   "true",
			CreatedAt:  now,
			CreatedBy:  "test-user",
			DeletedAt:  "",
		},
	}
}