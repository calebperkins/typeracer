package gameserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calebperkins.com/typeracer/internal/game"
)

func TestListGamesEmpty(t *testing.T) {
	gs := New()
	rec := httptest.NewRecorder()
	gs.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/games", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var got []game.Game
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("got %v, want empty non-nil list", got)
	}
}

func TestListGamesWithGames(t *testing.T) {
	gs := New()
	gs.games = append(gs.games, &game.Game{Id: "a"}, &game.Game{Id: "b"})
	rec := httptest.NewRecorder()
	gs.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/games", nil))

	var got []game.Game
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 || got[0].Id != "a" || got[1].Id != "b" {
		t.Errorf("got %+v, want games a and b", got)
	}
}

func TestListGamesMethodNotAllowed(t *testing.T) {
	gs := New()
	rec := httptest.NewRecorder()
	gs.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/games", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
