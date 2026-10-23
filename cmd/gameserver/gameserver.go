package main

import (
	"net/http"

	"calebperkins.com/typeracer/internal/api"
	"calebperkins.com/typeracer/internal/game"
)

func main() {
	server := api.New([]*game.Game{
		{Id: "foo"},
		{Id: "bar"},
	})
	http.ListenAndServe(":8090", server.Routes())
}
