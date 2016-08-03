package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Dummy variables used for testing
var p1 = Player{Name: "blueberry", ID: "muffin", Mark: O}
var p2 = Player{Name: "omelette", ID: "bacon", Mark: X}
var dummyBoard = map[string]string{A1: X, B1: empty, C1: O, A2: X, B2: empty, C2: O, A3: empty, B3: X, C3: empty}
var dummyGame = Game{Board: dummyBoard, Player1: p1, Player2: p2, CurrentPlayer: p2}
var dummyUsers = map[string]string{"blueberry": "muffin", "omelette": "bacon"}
var dummyChannelUsers = map[string][]string{"fakeChan": []string{"muffin", "bacon"}}

// Test represents one test (for any of the command functions).
type Test struct {
	name                 string
	existingBoards       map[string]*Game
	existingUsers        map[string]string
	existingChannelUsers map[string][]string
	input                string
	expectedResponse     string
}

// RunTest takes a list of Tests and a function to test them with and ensures that the response from running that
// function matches what the Test specifies is expected. It starts up a server to mock out the GameHandler, sets up
// any test environment variables (existingBoards, existingUsers, etc), makes a dummy request, and then compares
// the response to expectedResponse.
func RunTest(tests []Test, testFunc func([]string, RequestData) (*ResponseData, error), t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Header.Get("Input")
		inputList := strings.Split(input, " ")
		resp, err := testFunc(inputList, RequestData{channel: "fakeChan", username: "omelette", userID: "bacon"})
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		j, err := resp.getJSON()
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprint(w, string(j))
		return
	}))
	defer s.Close()

	for _, test := range tests {
		if test.existingBoards != nil {
			CurrentGames = test.existingBoards
		}
		if test.existingUsers != nil {
			Users = test.existingUsers
		}
		if test.existingChannelUsers != nil {
			ChannelUsers = test.existingChannelUsers
		}

		req, err := http.NewRequest("POST", s.URL, nil)
		if err != nil {
			t.Fatalf("Failed to create POST request: %v", err)
		}
		req.Header.Add("Input", test.input)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to do POST request: %v", err)
		}
		body, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read out response body: %v", err)
		}
		if test.expectedResponse != "" {
			if string(body) != test.expectedResponse {
				t.Errorf("Expected %s, but got %s when testing '%s'", test.expectedResponse, string(body), test.name)
			}
		}
		cleanup()
	}
}

// cleanup is used to reset all values that may have been changed during a test
func cleanup() {
	CurrentGames = map[string]*Game{}
	Users = map[string]string{}
	ChannelUsers = map[string][]string{}
	dummyGame.Board = map[string]string{A1: X, B1: empty, C1: O, A2: X, B2: empty, C2: O, A3: empty, B3: X, C3: empty}
	dummyGame.CurrentPlayer = p2
}

func TestStart(t *testing.T) {
	var successResponse = `{"response_type":"in_channel","text":"Hey \u003c@muffin|blueberry\u003e, omelette has challenged you to a game! Your move. \n      1    2    3\nA  ... | ... | ...\nB  ... | ... | ...\nC  ... | ... | ..."}`

	tests := []Test{
		Test{
			name:                 "basic",
			input:                "start @blueberry",
			existingUsers:        dummyUsers,
			existingChannelUsers: dummyChannelUsers,
			expectedResponse:     successResponse,
		},
		Test{
			name:             "no user",
			input:            "start",
			expectedResponse: UsageError.Error(),
		},
		Test{
			name:             "user doesn't exist",
			input:            "start @bacon",
			expectedResponse: UserDoesntExistError.Error(),
		},
		Test{
			name:                 "user in diff channel",
			input:                "start @blueberry",
			existingUsers:        dummyUsers,
			existingChannelUsers: map[string][]string{"fakeChan": []string{"bacon"}},
			expectedResponse:     UserNotInChannelError.Error(),
		},
		Test{
			name:                 "too many args",
			input:                "start @blueberry @pancakes",
			existingUsers:        dummyUsers,
			existingChannelUsers: dummyChannelUsers,
			expectedResponse:     UsageError.Error(),
		},
		Test{
			name:                 "bad username",
			input:                "start blueberry",
			existingUsers:        dummyUsers,
			existingChannelUsers: dummyChannelUsers,
			expectedResponse:     UsageError.Error(),
		},
		Test{
			name:                 "game already exists",
			input:                "start @blueberry",
			existingUsers:        dummyUsers,
			existingChannelUsers: dummyChannelUsers,
			existingBoards:       map[string]*Game{"fakeChan": &Game{}},
			expectedResponse:     GameAlreadyExistsError.Error(),
		},
		Test{
			name:                 "game exists in different channel",
			input:                "start @blueberry",
			existingUsers:        dummyUsers,
			existingChannelUsers: dummyChannelUsers,
			existingBoards:       map[string]*Game{"waffles": &Game{}},
			expectedResponse:     successResponse,
		},
	}
	RunTest(tests, handleStart, t)
}

func TestMove(t *testing.T) {
	var successResponse = `{"response_type":"in_channel","text":"blueberry (O) vs. omelette (X)\nX | X | ...\nX | ... | X\n0 | 0 | ...\nIt's blueberry's turn to make a move."}`
	tests := []Test{
		Test{
			name:          "move properly",
			input:         "move B1",
			existingUsers: dummyUsers,
			existingBoards: map[string]*Game{
				"fakeChan": &dummyGame,
			},
			expectedResponse: successResponse,
		},
		Test{
			name:             "no game exists",
			input:            "move B1",
			existingUsers:    dummyUsers,
			existingBoards:   map[string]*Game{"waffles": &Game{}},
			expectedResponse: NoGameExistsError.Error(),
		},
		Test{
			name:          "invalid move",
			input:         "move X1",
			existingUsers: dummyUsers,
			existingBoards: map[string]*Game{
				"fakeChan": &dummyGame,
			},
			expectedResponse: InvalidMoveError.Error(),
		},
		Test{
			name:          "position taken",
			input:         "move A1",
			existingUsers: dummyUsers,
			existingBoards: map[string]*Game{
				"fakeChan": &dummyGame,
			},
			expectedResponse: PositionTakenError.Error(),
		},
		Test{
			name:             "not authorized",
			input:            "move B1",
			existingUsers:    map[string]string{"blueberry": "muffin"},
			existingBoards:   map[string]*Game{"fakeChan": &Game{}},
			expectedResponse: NotAuthorizedError.Error(),
		},
		Test{
			name:             "too few args",
			input:            "move",
			existingBoards:   map[string]*Game{"fakeChan": &Game{}},
			expectedResponse: UsageError.Error(),
		},
		Test{
			name:             "too many args",
			input:            "move all the way!",
			existingBoards:   map[string]*Game{"fakeChan": &Game{}},
			expectedResponse: UsageError.Error(),
		},
		Test{
			name:          "not your turn",
			input:         "move B1",
			existingUsers: dummyUsers,
			existingBoards: map[string]*Game{
				"fakeChan": &Game{
					Board:         dummyBoard,
					Player1:       p1,
					Player2:       p2,
					CurrentPlayer: p1,
				},
			},
			expectedResponse: NotYourTurnError.Error(),
		},
	}
	RunTest(tests, handleMove, t)
}

func TestMultipleMoves(t *testing.T) {
	CurrentGames["fakeChan"] = &dummyGame
	req1 := RequestData{channel: "fakeChan", username: "omelette", userID: "bacon"}
	_, err := handleMove([]string{"move", "C3"}, req1)
	if err != nil {
		t.Fatalf("Failed to make valid move: %v", err)
	}
	req2 := RequestData{channel: "fakeChan", username: "blueberry", userID: "muffin"}
	_, err = handleMove([]string{"move", "A3"}, req2)
	if err != nil {
		t.Fatalf("Failed to make valid move: %v", err)
	}
	_, err = handleMove([]string{"move", "B1"}, req1)
	if err != nil {
		t.Fatalf("Failed to make valid move: %v", err)
	}
	_, err = handleMove([]string{"move", "A3"}, req2)
	if err == nil || err != PositionTakenError {
		t.Fatalf("Expected invalid move error, but got %v", err)
	}
	cleanup()
}

func TestGameEndsWithWinningMove(t *testing.T) {
	dummyGame.Board = map[string]string{A1: empty, B1: empty, C1: O, A2: X, B2: empty, C2: O, A3: X, B3: empty, C3: empty}
	CurrentGames["fakeChan"] = &dummyGame
	req1 := RequestData{channel: "fakeChan", username: "omelette", userID: "bacon"}
	_, err := handleMove([]string{"move", "A1"}, req1)
	if err != nil {
		t.Fatalf("Failed to make valid move: %v", err)
	}
	req2 := RequestData{channel: "fakeChan", username: "blueberry", userID: "muffin"}
	_, err = handleMove([]string{"move", "B1"}, req2)
	if err == nil || err != NoGameExistsError {
		t.Errorf("Game did not end when it should have")
	}
	cleanup()
}

func TestGameEndsWithTyingMove(t *testing.T) {
	dummyGame.Board = map[string]string{A1: empty, B1: X, C1: O, A2: X, B2: O, C2: O, A3: X, B3: O, C3: X}
	CurrentGames["fakeChan"] = &dummyGame
	req1 := RequestData{channel: "fakeChan", username: "omelette", userID: "bacon"}
	_, err := handleMove([]string{"move", "A1"}, req1)
	if err != nil {
		t.Fatalf("Failed to make valid move: %v", err)
	}
	req2 := RequestData{channel: "fakeChan", username: "blueberry", userID: "muffin"}
	_, err = handleMove([]string{"move", "B1"}, req2)
	if err == nil || err != NoGameExistsError {
		t.Errorf("Game did not end when it should have")
	}
	cleanup()
}

func TestDisplay(t *testing.T) {
	var successResponse = `{"response_type":"in_channel","text":"blueberry (O) vs. omelette (X)\nX | X | ...\n... | ... | X\n0 | 0 | ...\nIt's omelette's turn to make a move."}`
	tests := []Test{
		Test{
			name:          "display properly",
			input:         "display",
			existingUsers: dummyUsers,
			existingBoards: map[string]*Game{
				"fakeChan": &dummyGame,
			},
			expectedResponse: successResponse,
		},
		Test{
			name:             "no game exists",
			input:            "display",
			existingUsers:    dummyUsers,
			existingBoards:   map[string]*Game{"waffles": &Game{}},
			expectedResponse: NoGameExistsError.Error(),
		},
		Test{
			name:             "too many args",
			input:            "display everything!",
			existingBoards:   map[string]*Game{"fakeChan": &Game{}},
			expectedResponse: UsageError.Error(),
		},
	}
	RunTest(tests, handleDisplay, t)
}

func TestCancel(t *testing.T) {
	var successResponse = `{"response_type":"in_channel","text":"omelette has cancelled the current game. Perhaps a rematch later."}`
	tests := []Test{
		Test{
			name:             "cancel properly",
			input:            "cancel",
			existingBoards:   map[string]*Game{"fakeChan": &Game{Player1: Player{Name: "omelette", ID: "bacon"}}},
			expectedResponse: successResponse,
		},
		Test{
			name:             "no game exists",
			input:            "cancel",
			existingBoards:   map[string]*Game{"waffles": &Game{}},
			expectedResponse: NoGameExistsError.Error(),
		},
		Test{
			name:             "too many args",
			input:            "cancel everything!",
			existingBoards:   map[string]*Game{"fakeChan": &Game{}},
			expectedResponse: UsageError.Error(),
		},
		Test{
			name:             "not authorized",
			input:            "cancel",
			existingBoards:   map[string]*Game{"fakeChan": &Game{}},
			expectedResponse: NotAuthorizedError.Error(),
		},
	}
	RunTest(tests, handleCancel, t)
}
