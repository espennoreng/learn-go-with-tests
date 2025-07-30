package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
	"github.com/espennoreng/learn-go-with-tests/velo/pkg/api"
	"github.com/espennoreng/learn-go-with-tests/velo/testutils"
)

func TestHandlerRouting(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"get item", http.MethodGet, "/items/item-001", http.StatusOK},
		{"get nonexistent item", http.MethodGet, "/items/does-not-exist", http.StatusNotFound},
		{"get all items", http.MethodGet, "/items", http.StatusOK},
		{"update item", http.MethodPatch, "/items/item-001", http.StatusOK},
		{"delete item", http.MethodDelete, "/items/item-001", http.StatusNoContent},
		{"invalid method on items", http.MethodPost, "/items", http.StatusMethodNotAllowed},
		{"invalid method on item", http.MethodPost, "/items/item-001", http.StatusMethodNotAllowed},
		{"invalid path", http.MethodGet, "/unknown", http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := testutils.NewStubAppStore()
			store.Items = testutils.CreateTestItems()

			handler := api.NewHandler(store)

			var body []byte
			if tc.method == http.MethodPatch {
				body = []byte(`{"Name":"Updated Name"}`)
			}

			response := testutils.MakeRequest(t, handler, tc.method, tc.path, body)
			testutils.AssertStatus(t, response.Code, tc.expectedStatus)
		})
	}
}

func TestInternalServerErrors(t *testing.T) {
	t.Run("returns 500 when GetItems fails", func(t *testing.T) {
		// Create base store and error wrapper
		baseStore := testutils.NewStubAppStore()

		errorStore := &testutils.ErrorStore{
			AppStore:              baseStore,
			ShouldErrorOnGetItems: true,
		}

		handler := api.NewHandler(errorStore)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/items", nil)

		testutils.AssertStatus(t, response.Code, http.StatusInternalServerError)
	})
}

func TestHandlerContentType(t *testing.T) {
	t.Run("returns JSON content type", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/items/item-001", nil)

		testutils.AssertStatus(t, response.Code, http.StatusOK)
		testutils.AssertContentType(t, response, api.ContentTypeJSON)

	})
}

func TestGetItem(t *testing.T) {
	t.Run("returns item by id", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/items/item-001", nil)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var retrievedItem models.Item
		json.NewDecoder(response.Body).Decode(&retrievedItem)

		if retrievedItem.ID != "item-001" {
			t.Errorf("Expected item ID item-001, got %s", retrievedItem.ID)
		}
	})

}

func TestGetItems(t *testing.T) {
	t.Run("returns all items", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/items", nil)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var retrievedItems []models.Item
		json.NewDecoder(response.Body).Decode(&retrievedItems)

		testutils.AssertItemsContain(t, retrievedItems, "item-001", "item-002")

	})
}

func TestUpdateItem(t *testing.T) {
	t.Run("updates item fields", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		// First get original item
		getResponse := testutils.MakeRequest(t, handler, http.MethodGet, "/items/item-001", nil)
		var originalItem models.Item
		json.NewDecoder(getResponse.Body).Decode(&originalItem)

		// Update the item
		updateData := map[string]string{
			"Name": "Updated Name",
		}
		updateJSON, _ := json.Marshal(updateData)

		updateResponse := testutils.MakeRequest(t, handler, http.MethodPatch, "/items/item-001", updateJSON)
		testutils.AssertStatus(t, updateResponse.Code, http.StatusOK)

		var updatedItem models.Item
		json.NewDecoder(updateResponse.Body).Decode(&updatedItem)

		// Check updated field
		if updatedItem.Name != "Updated Name" {
			t.Errorf("Name not updated, got %q, want %q", updatedItem.Name, "Updated Name")
		}

		// Check fields that shouldn't change
		if updatedItem.ID != originalItem.ID {
			t.Errorf("ID changed unexpectedly, got %q, want %q", updatedItem.ID, originalItem.ID)
		}

		if updatedItem.ExternalID != originalItem.ExternalID {
			t.Errorf("ExternalID changed unexpectedly, got %q want %q", updatedItem.ExternalID, originalItem.ExternalID)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		invalidJSON := []byte(`{"Name": not-valid-json}`)

		response := testutils.MakeRequest(t, handler, http.MethodPatch, "/items/item-001", invalidJSON)
		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns 404 for non-existent item", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		updateData := map[string]string{
			"Name": "Updated Name",
		}
		updateJSON, _ := json.Marshal(updateData)

		response := testutils.MakeRequest(t, handler, http.MethodPatch, "/items/does-not-exist", updateJSON)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestDeleteItem(t *testing.T) {
	t.Run("deletes existing item", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		// Delete the item
		deleteResponse := testutils.MakeRequest(t, handler, http.MethodDelete, "/items/item-001", nil)
		testutils.AssertStatus(t, deleteResponse.Code, http.StatusNoContent)

		// Verify the item is gone
		getResponse := testutils.MakeRequest(t, handler, http.MethodGet, "/items/item-001", nil)
		testutils.AssertStatus(t, getResponse.Code, http.StatusNotFound)
	})

	t.Run("returns 404 for non-existent item", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Items = testutils.CreateTestItems()

		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodDelete, "/items/does-not-exist", nil)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
