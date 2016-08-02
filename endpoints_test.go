package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct {
	name             string
	existingBoards   map[string]GameBoard
	input            string
	expectedResponse string
}

func RunTest(tests []Test, testFunc func([]string, RequestData) (*ResponseData, error), t *testing.T) {
	Users["clim"] = "myID"
	Users["omelette"] = "cheese"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Header.Get("Input")
		inputList := strings.Split(input, " ")
		resp, err := testFunc(inputList, RequestData{channel: "fakeChan", username: "omelette"})
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
		} else {
			CurrentGames = map[string]GameBoard{}
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
	}
}

func TestStart(t *testing.T) {

	var successResponse = `{"response_type":"in_channel","text":"\u003c@myID|clim\u003e, omelette has challenged you to a duel! To accept this noble challenge, make the first move.","attachments":null}`

	tests := []Test{
		Test{
			name:             "basic",
			input:            "start @clim",
			expectedResponse: successResponse,
		},
		Test{
			name:             "no user",
			input:            "start",
			expectedResponse: UsageError.Error(),
		},
		Test{
			name:             "too many args",
			input:            "start @clim @pancakes",
			expectedResponse: UsageError.Error(),
		},

		Test{
			name:             "user doesn't exist",
			input:            "start @bacon",
			expectedResponse: UserDoesntExistError.Error(),
		},
		Test{
			name:             "game already exists",
			input:            "start @clim",
			existingBoards:   map[string]GameBoard{"fakeChan": GameBoard{}},
			expectedResponse: GameAlreadyExistsError.Error(),
		},
		Test{
			name:             "game exists in different channel",
			input:            "start @clim",
			existingBoards:   map[string]GameBoard{"differentChannel": GameBoard{}},
			expectedResponse: successResponse,
		},
	}
	RunTest(tests, startGame, t)
}

func TestMove(t *testing.T) {
	// TODO
}

func TestDisplay(t *testing.T) {
	// TODO
}

func TestHelp(t *testing.T) {
	// TODO
}

func TestCancel(t *testing.T) {
	var successResponse = `{"response_type":"in_channel","text":"omelette has cancelled the current game. What a shame.","attachments":null}`
	tests := []Test{
		Test{
			name:             "cancel properly",
			input:            "cancel",
			existingBoards:   map[string]GameBoard{"fakeChan": GameBoard{}},
			expectedResponse: successResponse,
		},
		Test{
			name:             "no game exists",
			input:            "cancel",
			existingBoards:   map[string]GameBoard{"bacon": GameBoard{}},
			expectedResponse: NoGameExistsError.Error(),
		},
	}
	RunTest(tests, handleCancel, t)
}
