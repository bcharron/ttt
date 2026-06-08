package main

import (
	"github.com/bcharron/ttt/internal/game"
)

func main() {
	game := game.NewGame(3)
	game.Draw()
}
