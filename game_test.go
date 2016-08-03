package main

import "testing"

func TestBoardDisplay(t *testing.T) {
	g := New(Player{Name: "batman", ID: "1", Mark: O}, Player{Name: "superman", ID: "2", Mark: X})
	g.Board = map[string]string{
		A1: empty,
		B1: X,
		C1: O,
		A2: empty,
		B2: empty,
		C2: X,
		A3: empty,
		B3: O,
		C3: X,
	}
	expected := `batman (O) vs. superman (X)
... | ... | ...
X | ... | 0
0 | X | X`
	d := g.Display()
	if d != expected {
		t.Errorf("Got \n%s\n but expected \n%s\n", d, expected)
	}
}

var winningBoards = []map[string]string{
	map[string]string{A1: X, B1: empty, C1: O, A2: X, B2: empty, C2: O, A3: X, B3: empty, C3: empty},
	map[string]string{A1: X, B1: X, C1: X, A2: O, B2: empty, C2: empty, A3: empty, B3: O, C3: empty},
	map[string]string{A1: O, B1: empty, C1: X, A2: empty, B2: O, C2: empty, A3: X, B3: empty, C3: O},
}

func TestGameWinner(t *testing.T) {
	g := New(Player{Name: "batman", ID: "1", Mark: O}, Player{Name: "superman", ID: "2", Mark: X})
	for _, board := range winningBoards {
		g.Board = board
		if !g.HasWinner() {
			t.Errorf("Game is won, but wasn't detected: %s", board)
		}
	}
}

func TestGameTied(t *testing.T) {
	g := New(Player{Name: "batman", ID: "1", Mark: O}, Player{Name: "superman", ID: "2", Mark: X})
	g.Board = map[string]string{A1: O, B1: X, C1: X, A2: O, B2: O, C2: X, A3: X, B3: O, C3: X}
	if !g.IsOver() {
		t.Errorf("Game is tied, but wasn't detected")
	}

}
