package main

import (
	"fmt"
	"strings"
)

// Game represents one instance of a tic tac toe game. The current layout of the board
// is stored in Board. CurrentPlayer keeps track of whose turn it is ("p1" or "p2").
type Game struct {
	Board         map[string]string // Current layout. Maps positions (e.g. A3) to values (X, O, or empty)
	CurrentPlayer Player            // Either Player1 or Player2
	Player1       Player            // Challenger, goes second. Mark = O
	Player2       Player            // Challengee, goes first. Mark = X
}

// Player represents a player of the game.
type Player struct {
	Name string // username
	ID   string // userID
	Mark string // X or O
}

// New returns a new game being played by p1 and p2 (where p1 invited p2 to play). p2 makes the first move.
func New(p1, p2 Player) *Game {
	game := Game{
		Board:         emptyBoard,
		CurrentPlayer: p2,
		Player1:       p1,
		Player2:       p2,
	}
	return &game
}

// Display prints out a visual representation of the game's current board configuration.
func (g *Game) Display() string {
	display := ""
	for i, pos := range boardPositions {
		if i%3 == 0 {
			display = fmt.Sprintf("%s\n%s", display, g.Board[pos])
		} else {
			display = fmt.Sprintf("%s | %s", display, g.Board[pos])
		}
	}
	return display
}

// HasWinner checks all 8 possible winning configurations of the board
// and returns true if the current board has a winner.
func (g *Game) HasWinner() bool {
	a1 := g.Board[A1]
	b1 := g.Board[B1]
	c1 := g.Board[C1]
	a2 := g.Board[A2]
	b2 := g.Board[B2]
	c2 := g.Board[C2]
	a3 := g.Board[A3]
	b3 := g.Board[B3]
	c3 := g.Board[C3]
	if a1 != empty && a1 == a2 && a2 == a3 ||
		b1 != empty && b1 == b2 && b2 == b3 ||
		c1 != empty && c1 == c2 && c2 == c3 ||
		a1 != empty && a1 == b1 && b1 == c1 ||
		a2 != empty && a2 == b2 && b2 == c2 ||
		a3 != empty && a3 == b3 && b3 == c3 ||
		a1 != empty && a1 == b2 && b2 == c3 ||
		a3 != empty && a3 == b2 && b2 == c1 {
		return true
	}
	return false
}

// HasTie checks if all the spots are nonempty (i.e. there's been a tie, since
// it is always called after checking HasWinner), and if so, returns true.
func (g *Game) HasTie() bool {
	if g.Board[A1] != empty &&
		g.Board[B1] != empty &&
		g.Board[C1] != empty &&
		g.Board[A2] != empty &&
		g.Board[B2] != empty &&
		g.Board[C2] != empty &&
		g.Board[A3] != empty &&
		g.Board[B3] != empty &&
		g.Board[C3] != empty {
		return true
	}
	return false
}

// moveIsValid is a helper function that returns true if a move is one of the
// 9 valid spots on the board.
func moveIsValid(move string) bool {
	for _, pos := range boardPositions {
		if strings.ToUpper(move) == pos {
			return true
		}
	}
	return false
}
