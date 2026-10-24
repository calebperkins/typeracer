package api

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"calebperkins.com/typeracer/internal/game"
	"github.com/coder/websocket"
)

func (gs *Server) getGameWebsocket(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("id")
	g := gs.findGameById(gameId)
	if g == nil {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	tokenCookie, err := r.Cookie(fmt.Sprintf("game-%s-token", gameId))
	var _ *game.Player
	if err == nil {
		_ = g.PlayerByToken(tokenCookie.Value)
	} else {
		// spectator
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("accept failed: %v", err)
		return
	}
	defer conn.CloseNow()

	updates, unsubscribe := g.Subscribe()
	defer unsubscribe()
	connCtx, stop := context.WithCancel(r.Context())
	defer stop()

	// Push every state change to this client. Writes are safe to run
	// concurrently with the reply writes below.
	go func() {
		for {
			select {
			case snap, ok := <-updates:
				if !ok {
					return
				}
				ctx, cancel := context.WithTimeout(connCtx, 10*time.Second)
				err := conn.Write(ctx, websocket.MessageText, wsReply(snap))
				cancel()
				if err != nil {
					conn.CloseNow() // unblocks the read loop
					return
				}
			case <-connCtx.Done():
				return
			}
		}
	}()

	// Block until the client disconnects; the deferred calls clean up after.
	for {
		// TODO send some player events
		if _, _, err := conn.Read(connCtx); err != nil {
			return // context cancelled or connection closed
		}
	}

}

type Message struct {
	Data string `json:"data"`
}

func wsReply(msg game.Snapshot) []byte {
	out, _ := json.Marshal(Message{Data: strconv.Itoa(int(msg.Version))})
	return out
}
