package main

import (
	"net/http"

	"calebperkins.com/typeracer/internal/gameserver"
)

func main() {
	gs := gameserver.New()
	http.ListenAndServe(":8090", gs.Routes())
}
