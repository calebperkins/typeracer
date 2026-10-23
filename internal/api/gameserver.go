package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"calebperkins.com/typeracer/internal/game"
)

// A Server hosts a set of games.
type Server struct {
	games []*game.Game
}

// New returns a Server hosting the given games.
func New(games []*game.Game) *Server {
	if games == nil {
		games = []*game.Game{}
	}
	return &Server{games: games}
}

func (gs *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /games", gs.listGames)
	mux.HandleFunc("GET /games/{id}", gs.getGame)
	return mux
}

func (gs *Server) listGames(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, gs.games)
}

func (gs *Server) getGame(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	g := gs.findGameById(id)
	if g == nil {
		http.Error(w, fmt.Sprintf("game %q not found", id), http.StatusNotFound)
		return
	}
	writeJSON(w, g)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (gs *Server) findGameById(id string) *game.Game {
	for _, g := range gs.games {
		if g.Id == id {
			return g
		}
	}
	return nil
}
