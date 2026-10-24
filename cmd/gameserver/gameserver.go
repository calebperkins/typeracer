package main

import (
	"net/http"

	"calebperkins.com/typeracer/internal/api"
	"calebperkins.com/typeracer/internal/game"
)

func main() {
	server := api.New([]*game.Game{
		game.New("foo", "bear"),
		game.New("bar", "bear"),
	})
	http.ListenAndServe(":8090", server.Routes())
}
