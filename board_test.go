package main

import "testing"

func TestBoardDisplay(t *testing.T) {
	b := New(Player{Name: "batman", ID: "1"}, Player{Name: "superman", ID: "2"})
	b.Board = map[string]string{
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
 X  | ... |  O 
 O  |  X  |  X 
It's superman's turn to make a move.`
	d := b.Display()
	if d != expected {
		t.Errorf("Got \n%s\n but expected \n%s\n", d, expected)
	}
}
