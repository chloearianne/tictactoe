package main

import "strings"

// GameBoard represents one instance of a tic tac toe game. The current layout of the board
// is stored in CurrentConfig, and CurrentPlayer keeps track of whose turn it is.
type GameBoard struct {
	CurrentConfig map[string]string
	CurrentPlayer string
	Players       []string
	GameOver      bool // FIXME: unnecessary?
}

var validPositions = []string{"A1", "A2", "A3", "B1", "B2", "B3", "C1", "C2", "C3"}

func moveIsValid(move string) bool {
	for _, pos := range validPositions {
		if strings.ToUpper(move) == pos {
			return true
		}
	}
	return false
}
