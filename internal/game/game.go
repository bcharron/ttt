package game

import (
	"errors"
	"fmt"
)

var ErrInvalidPosition = errors.New("invalid position")
var ErrAlreadyOccupied = errors.New("already occupied")
var ErrInvalidPlayer = errors.New("invalid player")

const (
	PLAYER_X    = -1
	EMPTY_SPACE = 0
	PLAYER_O    = 1
)

type Player int

type Game struct {
	grid      []int
	rows      int
	cols      int
	size      int // Total size
	movesLeft int // How many moves are left
}

func NewGame(rows int) *Game {
	size := rows * rows
	return &Game{
		grid:      make([]int, size),
		rows:      rows,
		cols:      rows, // it's square
		size:      rows * rows,
		movesLeft: size,
	}
}

func (g *Game) GetState() []int {
	out := make([]int, len(g.grid))
	copy(out, g.grid)
	return out
}

func Itoc(v int) string {
	switch v {
	case PLAYER_X:
		return "x"
	case EMPTY_SPACE:
		return " "
	case PLAYER_O:
		return "o"
	}

	return " "
}

func (g *Game) Draw() {
	fmt.Printf("+---+---+---+\n")
	fmt.Printf("+ %s | %s | %s |\n", Itoc(g.grid[0]), Itoc(g.grid[1]), Itoc(g.grid[2]))
	fmt.Printf("+---+---+---+\n")
	fmt.Printf("+ %s | %s | %s |\n", Itoc(g.grid[3]), Itoc(g.grid[4]), Itoc(g.grid[5]))
	fmt.Printf("+---+---+---+\n")
	fmt.Printf("+ %s | %s | %s |\n", Itoc(g.grid[6]), Itoc(g.grid[7]), Itoc(g.grid[8]))
	fmt.Printf("+---+---+---+\n")
}

func (g *Game) IsLegalMove(player int, pos int) error {
	if pos < 0 || pos > len(g.grid)-1 {
		return ErrInvalidPosition
	}

	if g.grid[pos] != EMPTY_SPACE {
		return ErrAlreadyOccupied
	}

	if player != PLAYER_X && player != PLAYER_O {
		return ErrInvalidPlayer
	}

	return nil
}

func (g *Game) Play(player int, pos int) error {
	if err := g.IsLegalMove(player, pos); err != nil {
		return err
	}

	g.grid[pos] = player

	g.movesLeft--

	return nil
}

func (g *Game) HaveMovesLeft() bool {
	return g.movesLeft > 0
}

func (g *Game) HaveWinner() (bool, int) {
	for row := 0; row < g.rows; row++ {
		if haveWinner, winner := g.horizontalWinner(row); haveWinner {
			return true, winner
		}
	}

	for col := 0; col < g.cols; col++ {
		if haveWinner, winner := g.verticalWinner(col); haveWinner {
			return true, winner
		}
	}

	if haveWinner, winner := g.diagonalWinner(); haveWinner {
		return true, winner
	}

	if haveWinner, winner := g.reverseDiagonalWinner(); haveWinner {
		return true, winner
	}

	return false, 0
}

func isPlayer(p int) bool {
	if p == PLAYER_X || p == PLAYER_O {
		return true
	}

	return false
}

func (g *Game) horizontalWinner(row int) (bool, int) {
	if row < 0 || row > g.rows-1 {
		panic("invalid row")
	}

	winner := g.get(row, 0)

	if !isPlayer(winner) {
		return false, 0
	}

	for col := 0; col < g.cols; col++ {
		if winner != g.get(row, col) {
			return false, 0
		}
	}

	return true, winner
}

func (g *Game) verticalWinner(col int) (bool, int) {
	if col < 0 || col > g.cols-1 {
		panic("invalid column")
	}

	winner := g.get(0, col)

	if !isPlayer(winner) {
		return false, 0
	}

	for row := 0; row < g.rows; row++ {
		if winner != g.get(row, col) {
			return false, 0
		}
	}

	return true, winner
}

func (g *Game) diagonalWinner() (bool, int) {
	winner := g.get(0, 0)

	if !isPlayer(winner) {
		return false, 0
	}

	for x := 0; x < g.rows; x++ {
		if g.get(x, x) != winner {
			return false, 0
		}
	}

	return true, winner
}

func (g *Game) reverseDiagonalWinner() (bool, int) {
	winner := g.get(0, g.cols-1)

	if !isPlayer(winner) {
		return false, 0
	}

	col := g.cols - 1
	for row := 0; row < g.rows; row++ {
		if winner != g.get(row, col) {
			return false, 0
		}

		col--
	}

	return true, winner
}

func (g *Game) get(row, col int) int {
	return g.grid[g.cols*row+col]
}
