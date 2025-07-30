package velo_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/espennoreng/learn-go-with-tests/velo/models"
	"github.com/espennoreng/learn-go-with-tests/velo/testutils"
)

func TestGETItems(t *testing.T) {

	t.Run("returns single item", func(t *testing.T) {
		server := testutils.SetupTestServer()

		id := "item-001"

		response := testutils.MakeRequest(t, server, http.MethodGet, fmt.Sprintf("/items/%s", id), nil)

		var retrievedItem models.Item
		err := json.NewDecoder(response.Body).Decode(&retrievedItem)

		testutils.AssertValidJSON(t, response.Body, err)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns all items", func(t *testing.T) {
		server := testutils.SetupTestServer()

		response := testutils.MakeRequest(t, server, http.MethodGet, "/items", nil)

		var retrievedItems []models.Item
		err := json.NewDecoder(response.Body).Decode(&retrievedItems)

		testutils.AssertValidJSON(t, response.Body, err)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		testutils.AssertItemsContain(t, retrievedItems, "item-001", "item-002")
	})

	t.Run("returns 404 on missing item", func(t *testing.T) {
		server := testutils.SetupTestServer()

		response := testutils.MakeRequest(t, server, http.MethodGet, "/items/does-not-exist", nil)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns items as JSON", func(t *testing.T) {
		server := testutils.SetupTestServer()

		response := testutils.MakeRequest(t, server, http.MethodGet, "/items/item-001", nil)

		var retrievedItem models.Item
		err := json.NewDecoder(response.Body).Decode(&retrievedItem)

		testutils.AssertValidJSON(t, response.Body, err)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		testutils.AssertContentType(t, response, "application/json")
	})

	t.Run("delete item", func(t *testing.T) {
		server := testutils.SetupTestServer()

		id := "item-001"

		// Delete the item
		deleteResponse := testutils.MakeRequest(t, server, http.MethodDelete, fmt.Sprintf("/items/%s", id), nil)
		testutils.AssertStatus(t, deleteResponse.Code, http.StatusNoContent)

		// Verify item is gone
		getResponse := testutils.MakeRequest(t, server, http.MethodGet, fmt.Sprintf("/items/%s", id), nil)
		testutils.AssertStatus(t, getResponse.Code, http.StatusNotFound)
	})

	t.Run("update item", func(t *testing.T) {
		server := testutils.SetupTestServer()

		id := "item-001"

		// get original item
		getResponse := testutils.MakeRequest(t, server, http.MethodGet, fmt.Sprintf("/items/%s", id), nil)
		testutils.AssertStatus(t, getResponse.Code, http.StatusOK)

		var originalItem models.Item
		json.NewDecoder(getResponse.Body).Decode(&originalItem)

		newName := "Updated Test Item"
		// payload with modified name
		updateData := map[string]string{
			"Name": newName,
		}
		updateJson, err := json.Marshal(updateData)

		if err != nil {
			t.Fatalf("error parsing updateData to json")
		}

		// patch request
		patchResponse := testutils.MakeRequest(t, server, http.MethodPatch, fmt.Sprintf("/items/%s", id), updateJson)
		testutils.AssertStatus(t, patchResponse.Code, http.StatusOK)

		var updatedItem models.Item
		json.NewDecoder(patchResponse.Body).Decode(&updatedItem)

		// assert that new field is updated
		if updatedItem.Name != newName {
			t.Errorf("Name field not updated, got %q, want %q", updatedItem.Name, newName)
		}

		// assert that other fields remind the same
		if updatedItem.ID != originalItem.ID {
			t.Errorf("ID changed unexpectedly, got %q, want %q", updatedItem.ID, originalItem.OrgID)
		}

		if updatedItem.ExternalID != originalItem.ExternalID {
			t.Errorf("ExternalID changed unexpectedly, got %q want %q", updatedItem.ExternalID, originalItem.ExternalID)
		}
	})

	t.Run("handles concurrent updates correctly", func(t *testing.T) {
		spy, server := testutils.SetupSpyTestServer()
		id := "item-001"
		concurrentUpdates := 10000

		// Get initial version
		getResponse := testutils.MakeRequest(t, server, http.MethodGet, fmt.Sprintf("/items/%s", id), nil)
		var originalItem models.Item
		json.NewDecoder(getResponse.Body).Decode(&originalItem)

		// Use a wait group to coordinate goroutines
		var wg sync.WaitGroup
		successCount := atomic.Int32{}

		for i := range concurrentUpdates {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// Get current item first (to get latest version)
				getResp := testutils.MakeRequest(t, server, http.MethodGet, fmt.Sprintf("/items/%s", id), nil)
				var item models.Item
				json.NewDecoder(getResp.Body).Decode(&item)

				update := map[string]any{
					"Name": fmt.Sprintf("Update %d", index),
				}
				updateJson, _ := json.Marshal(update)

				resp := testutils.MakeRequest(t, server, http.MethodPatch, fmt.Sprintf("/items/%s", id), updateJson)
				if resp.Code == http.StatusOK {
					successCount.Add(1)
				}
			}(i)
		}

		wg.Wait()

		// Verify final state
		finalResponse := testutils.MakeRequest(t, server, http.MethodGet, fmt.Sprintf("/items/%s", id), nil)
		var finalItem models.Item
		json.NewDecoder(finalResponse.Body).Decode(&finalItem)

		// Log spy metrics
		t.Logf("GetItem calls: %d", spy.GetItemCalls.Load())
		t.Logf("UpdateItem calls: %d", spy.UpdateItemCalls.Load())
		t.Logf("Success count: %d", successCount.Load())

		// Verify expectations with spy data
		if spy.UpdateItemCalls.Load() != int32(concurrentUpdates) {
			t.Errorf("Expected %d update calls, got %d", concurrentUpdates, spy.UpdateItemCalls.Load())
		}

		if spy.UpdateItemCalls.Load() != successCount.Load() {
			t.Errorf("Not all updates succeeded: %d attempts, %d successes",
				spy.UpdateItemCalls.Load(), successCount.Load())
		}

		// Check the last update in the log
		if len(spy.UpdateLog) > 0 {
			t.Logf("Last update: %s", spy.UpdateLog[len(spy.UpdateLog)-1])
		}
	})

}
