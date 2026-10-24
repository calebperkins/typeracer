package api

import (
	"encoding/json/v2"
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
	mux.HandleFunc("POST /games/{id}", gs.refresh)
	mux.HandleFunc("GET /games/{id}/ws", gs.getGameWebsocket)
	mux.HandleFunc("POST /games/{id}/players", gs.joinGame)
	return mux
}

func (gs *Server) refresh(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	g := gs.findGameById(id)
	if g == nil {
		http.Error(w, fmt.Sprintf("game %q not found", id), http.StatusNotFound)
		return
	}
	g.Refresh()
	writeJSON(w, g.Snapshot())
}

func (gs *Server) listGames(w http.ResponseWriter, r *http.Request) {
	snaps := make([]game.Snapshot, 0, len(gs.games))
	for _, g := range gs.games {
		snaps = append(snaps, g.Snapshot())
	}
	writeJSON(w, snaps)
}

func (gs *Server) getGame(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	g := gs.findGameById(id)
	if g == nil {
		http.Error(w, fmt.Sprintf("game %q not found", id), http.StatusNotFound)
		return
	}
	writeJSON(w, g.Snapshot())
}

type joinGameRequestBody struct {
	name string
}

func (gs *Server) joinGame(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	g := gs.findGameById(id)
	if g == nil {
		http.Error(w, fmt.Sprintf("game %q not found", id), http.StatusNotFound)
		return
	}
	var playerReq joinGameRequestBody
	if err := json.UnmarshalRead(r.Body, &playerReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := g.Join(playerReq.name)
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: p.Token,
	})
	writeJSON(w, p)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.MarshalWrite(w, v)
	if err != nil {
		panic(err)
	}
}

func (gs *Server) findGameById(id string) *game.Game {
	for _, g := range gs.games {
		if g.Id() == id {
			return g
		}
	}
	return nil
}
