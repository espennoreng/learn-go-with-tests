package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingTem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(t, player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest(t))
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player {
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
}
