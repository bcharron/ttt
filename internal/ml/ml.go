package ml

import (
	"errors"
	"fmt"
	"github.com/bcharron/mllib"
	"github.com/bcharron/ttt/internal/game"
	"math/rand"
	"sort"
)

var ErrNoLegalPlay = errors.New("no legal moves to play")

func GridToInputs(state []int, player int) []float64 {
	output := make([]float64, len(state))

	// Present the grid in a way that "1" is always "this player", and "-1"
	// is the opponent.
	for x, v := range state {
		switch v {
		case game.EMPTY_SPACE:
			output[x] = 0
		case player:
			output[x] = 1
		default:
			output[x] = -1
		}
	}

	return output
}

func sortArrayByIndexDescending(array []float64) []int {
	indices := make([]int, len(array))
	for x := range indices {
		indices[x] = x
	}

	sort.Slice(indices, func(i, j int) bool {
		return array[indices[i]] > array[indices[j]]
	})

	return indices
}

// epsilon is the probability that PlayNextMove will choose a random move
// rather than the "best" move. This is necessary to learn, otherwise X and O
// will stalemate over and over during training.
func PlayNextMove(player int, g *game.Game, n *mllib.Network, epsilon float64) (GameRecord, error) {
	inputs := GridToInputs(g.GetState(), player)
	predictions := make([]float64, len(inputs))
	n.Forward(predictions, inputs)

	sortedMoves := sortArrayByIndexDescending(predictions)

	// fmt.Printf("predictions: %v\n", predictions)

	preState := g.GetState()

	if rand.Float64() <= epsilon {
		rand.Shuffle(len(sortedMoves), func(i, j int) {
			sortedMoves[i], sortedMoves[j] = sortedMoves[j], sortedMoves[i]
		})
	}

	for _, pos := range sortedMoves {
		if g.Play(player, pos) == nil {
			record := GameRecord{
				Player: player,
				Move:   pos,
				State:  preState,
			}

			return record, nil
		}
	}

	return GameRecord{}, ErrNoLegalPlay
}

type GameRecord struct {
	Player int
	State  []int
	Move   int
}

func PlayGame(g *game.Game, nn []*mllib.Network) (int, []GameRecord) {
	turn := 0
	players := []int{game.PLAYER_X, game.PLAYER_O}
	history := make([]GameRecord, 0)

	for {
		player := players[turn%2]
		net := nn[turn%2]

		record, err := PlayNextMove(player, g, net, 0.1)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return 0, history
		}

		// record := GameRecord{
		// 	Player:      player,
		// 	State:       g.GetState(),
		// 	Predictions: predictions,
		// 	Move:        move,
		// }

		history = append(history, record)

		if ok, winner := g.HaveWinner(); ok {
			// fmt.Printf("Winner: %s\n", game.Itoc(winner))
			return winner, history
		}

		if !g.HaveMovesLeft() {
			// fmt.Printf("It's a draw!\n")
			return 0, history
		}

		turn++
	}

	// shouldn't happen
	return 0, history
}
