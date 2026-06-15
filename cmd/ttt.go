package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/bcharron/ttt/internal/game"
)

func getMove() int {
	fmt.Printf("Enter pos\n")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())

		if err == nil {
			return n
		} else {
			fmt.Printf("Position must be an integer\n")
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("input error: %s\n", scanner.Err().Error())
		os.Exit(1)
	}

	panic("uhh")
}

func main() {
	g := game.NewGame(3)

	turn := 0
	players := []int{game.PLAYER_X, game.PLAYER_O}

	for {
		player := players[turn%2]

		g.Draw()

		if ok, winner := g.HaveWinner(); ok {
			fmt.Printf("Winner: %s\n", game.Itoc(winner))
			break
		}

		if !g.HaveMovesLeft() {
			fmt.Printf("It's a draw!\n")
			break
		}

		for {
			fmt.Printf("Player %s's turn\n", game.Itoc(player))
			pos := getMove()

			err := g.Play(player, pos)

			if err == nil {
				break
			} else {
				fmt.Printf("Illegal move: %s\n", err.Error())
			}
		}

		turn++
	}
}
