package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/bcharron/mllib"
	"github.com/bcharron/ttt/internal/game"
	"github.com/bcharron/ttt/internal/ml"
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

func playerToIndex(player int) int {
	switch player {
	case game.PLAYER_X:
		return 0
	case game.PLAYER_O:
		return 1
	}

	return 2
}

func main() {
	const epochs = 5000

	nn := make([]*mllib.Network, 20)

	// player 1 (X)
	// 9 inputs, two hidden layers with 16 inputs, 9 outputs
	nn[0] = mllib.NewNetwork(mllib.Tanh, 9, 16, 16, 9)
	nn[0].RandomizeWeights()

	// player 2 (O)
	nn[1] = mllib.NewNetwork(mllib.Tanh, 9, 16, 16, 9)
	nn[1].RandomizeWeights()

	wins := make([]int, 3)

	// Train two networks against each other
	for range epochs {
		g := game.NewGame(3)
		winner, history := ml.PlayGame(g, nn)
		g.Draw()
		fmt.Printf("winner: %v\n", game.Itoc(winner))

		winnerIdx := playerToIndex(winner)

		wins[winnerIdx]++

		expected := make([]float64, 9)
		for turn, record := range history {
			var reward float64
			playerIndex := playerToIndex(record.Player)

			if record.Player == winner { // winner
				reward = 1
			} else if winner != 0 { // loser
				reward = -1
			} else {
				reward = 0
			}

			tmp := make([]float64, 9)

			// Must run Forward to compute the cache elements necessary for Backward
			nn[playerIndex].Forward(tmp, ml.GridToInputs(record.State, record.Player))

			for i := range record.State {
				if i == record.Move {
					expected[i] = tmp[i] + reward
				} else {
					expected[i] = tmp[i]
				}
			}

			fmt.Printf("turn %d, player %d, move = %d, reward = %.2f\n", turn, record.Player, record.Move, reward)

			nn[playerIndex].Backward(expected, 0.05)
		}
	}

	fmt.Printf("Wins: X=%d, O=%d, Ties=%d\n", wins[0], wins[1], wins[2])

	turn := 0
	players := []int{game.PLAYER_X, game.PLAYER_O}

	g := game.NewGame(3)
	g.Draw()

	var net *mllib.Network
	if wins[0] > wins[1] {
		fmt.Printf("Playing against X\n")
		net = nn[0]
	} else {
		fmt.Printf("Playing against O\n")
		net = nn[1]
	}

	// Play a game against human
	for {
		player := players[turn%2]
		if turn%2 == 0 {
			_, err := ml.PlayNextMove(player, g, net, 0.0)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				break
			}
		} else {
			fmt.Printf("Player %s's turn\n", game.Itoc(player))
			pos := getMove()

			err := g.Play(player, pos)

			if err != nil {
				fmt.Printf("Illegal move: %s\n", err.Error())
				break
			}
		}

		g.Draw()

		if ok, winner := g.HaveWinner(); ok {
			fmt.Printf("Winner: %s\n", game.Itoc(winner))
			break
		}

		if !g.HaveMovesLeft() {
			fmt.Printf("It's a draw!\n")
			break
		}

		turn++
	}
}
