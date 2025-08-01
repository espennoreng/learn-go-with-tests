package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
	"github.com/espennoreng/learn-go-with-tests/velo/pkg/api"
	"github.com/espennoreng/learn-go-with-tests/velo/testutils"
)

func TestInvalidRouting(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body 			[]byte
		expectedStatus int
	}{
		{"create item with empty body", http.MethodPost, "/items", nil, http.StatusBadRequest},
		{"invalid method on /items", http.MethodPatch, "/items", nil, http.StatusMethodNotAllowed},		
		{"invalid method on /users", http.MethodPatch, "/users", nil, http.StatusMethodNotAllowed},
		{"invalid method on /items/{}", http.MethodPost, "/items/item-001", nil, http.StatusMethodNotAllowed},
		{"invalid method on /users/{}", http.MethodDelete, "/users/random-id", nil, http.StatusMethodNotAllowed},
		{"invalid method on sessions", http.MethodPatch, "/sessions/session-001", nil, http.StatusMethodNotAllowed},
		{"invalid path", http.MethodGet, "/unknown", nil, http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := testutils.NewStubAppStore()
			store.Items = testutils.CreateTestItems()
			store.Sessions = testutils.CreateTestSessions()

			handler := api.NewHandler(store)

			response := testutils.MakeRequest(t, handler, tc.method, tc.path, tc.body)
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
			ShouldError: true,
		}

		handler := api.NewHandler(errorStore)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/items", nil)

		testutils.AssertStatus(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("returns 500 when CreateItem fails", func(t *testing.T) {
		// Create base store and error wrapper
		baseStore := testutils.NewStubAppStore()

		errorStore := &testutils.ErrorStore{
			AppStore:              baseStore,
			ShouldError: true,
		}

		handler := api.NewHandler(errorStore)

		createItemData := api.CreateItemRequest{Name: "a name"}
		body, err := json.Marshal(createItemData)
		if err != nil {
			t.Fatalf("could not marshal JSON: %v", err)
		}

		response := testutils.MakeRequest(t, handler, http.MethodPost, "/items", body)

		testutils.AssertStatus(t, response.Code, http.StatusInternalServerError)
	})

		t.Run("returns 500 when CreateUser fails", func(t *testing.T) {
		// Create base store and error wrapper
		baseStore := testutils.NewStubAppStore()

		errorStore := &testutils.ErrorStore{
			AppStore:              baseStore,
			ShouldError: true,
		}

		handler := api.NewHandler(errorStore)

		createUserData := api.CreateItemRequest{Name: "a name"}
		body, err := json.Marshal(createUserData)
		if err != nil {
			t.Fatalf("could not marshal JSON: %v", err)
		}

		response := testutils.MakeRequest(t, handler, http.MethodPost, "/users", body)

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

func TestCreateItem(t *testing.T){
	t.Run("create item", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		itemName := "newly created item"

		createItemData := api.CreateItemRequest{Name: itemName}
		body, err := json.Marshal(createItemData)
		if err != nil {
			t.Fatalf("could not marshal JSON: %v", err)
		}

		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/items", body)

		testutils.AssertStatus(t, createResponse.Code, http.StatusCreated)

		location := createResponse.Header().Get("Location")
		if location == ""{
			t.Fatalf("expected Location header to be set, but it was empty")
		}

		var createdItem models.Item
		json.NewDecoder(createResponse.Body).Decode(&createdItem)

		if createdItem.Name != itemName {
			t.Errorf("expected item name %q, got %q", itemName, createdItem.Name)
		}

		getResponse := testutils.MakeRequest(t, handler, http.MethodGet, fmt.Sprintf("/items/%s", createdItem.ID), nil)

		testutils.AssertStatus(t, getResponse.Code, http.StatusOK)
	})

	t.Run("create item with bad JSON", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		badJSON := []byte(`{"name": "New item}`)
		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/items", badJSON)

		testutils.AssertStatus(t, createResponse.Code, http.StatusBadRequest)
	})

	t.Run("create item with invalid data", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		invalidData := []byte(`{"name": ""}`)
		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/items", invalidData)

		testutils.AssertStatus(t, createResponse.Code, http.StatusBadRequest)
	})

	t.Run("create item with empty body", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/items", nil)

		testutils.AssertStatus(t, createResponse.Code, http.StatusBadRequest)
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

		testutils.AssertContainsID(t, retrievedItem, "item-001")
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

		testutils.AssertContainsIDs(t, retrievedItems, "item-001", "item-002")

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

func TestGetSession(t *testing.T){
	t.Run("returns item by id", func (t *testing.T)  {
		store := testutils.NewStubAppStore()
		store.Sessions = testutils.CreateTestSessions()

		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/sessions/session-001", nil)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var retrievedSession models.Session
		json.NewDecoder(response.Body).Decode(&retrievedSession)

		testutils.AssertContainsID(t, retrievedSession, "session-001")
	})

	t.Run("returns 404 for no-existent session", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		store.Sessions = testutils.CreateTestSessions()

		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/sessions/does-not-exist", nil)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}


func TestUsersHandler(t *testing.T){
	t.Run("GET /users/{id}", func(t *testing.T) {
		store := testutils.NewStubAppStoreWithData()
		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/users/user-001", nil)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var receivedUser models.User
		json.NewDecoder(response.Body).Decode(&receivedUser)
	
		testutils.AssertContainsID(t, receivedUser, "user-001")
	})

	t.Run("GET /users/{non-existing ID}", func(t *testing.T) {
		store := testutils.NewStubAppStoreWithData()
		handler := api.NewHandler(store)

		response := testutils.MakeRequest(t, handler, http.MethodGet, "/users/does-not-exist", nil)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("POST /users", func (t *testing.T)  {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		newUserName := "Per"
		createUserData := api.CreateUserRequest{Name: newUserName}
		body, err := json.Marshal(createUserData)
		if err != nil {
			t.Fatalf("could not marshal JSON: %v", err)
		}

		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/users", body)
		
		var createdUser models.User
		json.NewDecoder(createResponse.Body).Decode(&createdUser)

		if createUserData.Name != newUserName{
			t.Errorf("got %q, want %q", createUserData.Name, newUserName)
		}

		location := createResponse.Header().Get("Location")
		if location == ""{
			t.Fatalf("expected Location header to be set, but it was empty")
		}

		getResponse := testutils.MakeRequest(t, handler, http.MethodGet, fmt.Sprintf("/users/%s", createdUser.ID), nil)

		testutils.AssertStatus(t, getResponse.Code, http.StatusOK)
	})

	t.Run("POST /users with bad JSON", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		badJSON := []byte(`{"name": "Per}`)
		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/users", badJSON)

		testutils.AssertStatus(t, createResponse.Code, http.StatusBadRequest)
	})

	t.Run("POST /users with invalid data", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		invalidData := []byte(`{"name": ""}`)
		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/users", invalidData)

		testutils.AssertStatus(t, createResponse.Code, http.StatusBadRequest)
	})

	t.Run("POST /users with empty body", func(t *testing.T) {
		store := testutils.NewStubAppStore()
		handler := api.NewHandler(store)

		createResponse := testutils.MakeRequest(t, handler, http.MethodPost, "/users", nil)

		testutils.AssertStatus(t, createResponse.Code, http.StatusBadRequest)
	})
}