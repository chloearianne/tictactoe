package main

import (
	"fmt"
	"strings"
)

// Game represents one instance of a tic tac toe game. The current layout of the board
// is stored in Board. CurrentPlayer keeps track of whose turn it is ("p1" or "p2").
type Game struct {
	Board         map[string]string
	CurrentPlayer string // Either p1 or p2
	Player1       Player // Challenger, goes second. Mark = O
	Player2       Player // Challengee, goes first. Mark = X
}

// Player represents a player of the game.
type Player struct {
	Name string
	ID   string
}

// New returns a new game being played by p1 and p2, where p1 invited p2 to play. p2 makes the first move.
func New(p1, p2 Player) *Game {
	game := Game{
		Board:         emptyBoard,
		CurrentPlayer: "p2",
		Player1:       p1,
		Player2:       p2,
	}
	return &game
}

func (g *Game) Display() string {
	display := fmt.Sprintf("%s (O) vs. %s (X)", g.Player1.Name, g.Player2.Name)
	for i, pos := range boardPositions {
		if i%3 == 0 {
			display = fmt.Sprintf("%s\n%s", display, g.Board[pos])
		} else {
			display = fmt.Sprintf("%s | %s", display, g.Board[pos])
		}
	}
	var player string
	if g.CurrentPlayer == "p1" {
		player = g.Player1.Name
	} else {
		player = g.Player2.Name
	}
	display = fmt.Sprintf("%s\nIt's %s's turn to make a move.", display, player)
	return display
}

func moveIsValid(move string) bool {
	for _, pos := range boardPositions {
		if strings.ToUpper(move) == pos {
			return true
		}
	}
	return false
}
