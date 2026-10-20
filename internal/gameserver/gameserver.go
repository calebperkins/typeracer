package gameserver

import (
	"encoding/json"
	"net/http"

	"calebperkins.com/typeracer/internal/game"
)

// A GameServer hosts a set of games.
type GameServer struct {
	games []*game.Game
}

func New() *GameServer {
	return &GameServer{games: []*game.Game{}}
}

func (gs *GameServer) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /games", gs.listGames)
	return mux
}

func (gs *GameServer) listGames(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, gs.games)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
