package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidPositions(t *testing.T) {
	g := NewGame(3)

	assert.ErrorIs(t, ErrInvalidPlayer, g.Play(2, 1))

	// Position outside 0..8
	assert.ErrorIs(t, ErrInvalidPosition, g.Play(PLAYER_X, -1))
	assert.ErrorIs(t, ErrInvalidPosition, g.Play(PLAYER_X, 9))

	for x := 0; x < g.size; x++ {
		assert.Nil(t, g.Play(PLAYER_X, x))
	}

	assert.ErrorIs(t, ErrAlreadyOccupied, g.Play(PLAYER_X, 0))

}

func TestWinning(t *testing.T) {
	g := NewGame(3)

	// row 1
	for x := 0; x < g.size; x++ {
		g.grid[x] = PLAYER_X
	}

	haveWinner, winner := g.HaveWinner()
	assert.True(t, haveWinner)
	assert.Equal(t, PLAYER_X, winner)

	g = NewGame(3)
	haveWinner, winner = g.HaveWinner()
	assert.False(t, haveWinner)

	// Diagonal winner
	g.Play(PLAYER_X, 0)
	g.Play(PLAYER_X, 4)
	g.Play(PLAYER_X, 8)

	haveWinner, winner = g.HaveWinner()
	assert.True(t, haveWinner)
	assert.Equal(t, PLAYER_X, winner)

	g = NewGame(3)
	haveWinner, winner = g.HaveWinner()
	assert.False(t, haveWinner)

	// Reverse-Diagonal winner
	g.Play(PLAYER_X, 2)
	g.Play(PLAYER_X, 4)
	g.Play(PLAYER_X, 6)

	g.Draw()

	haveWinner, winner = g.HaveWinner()
	assert.True(t, haveWinner)
	assert.Equal(t, PLAYER_X, winner)
}
