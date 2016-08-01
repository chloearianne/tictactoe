package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var UsageError = errors.New(`Use /ttt to play a game of tic tac toe.
To start a game: /ttt start [@user]
To make a move: /ttt move [position]
To display current board: /ttt display
For help: /ttt help`)

// RequestData represents an incoming request from a Slack channel when a user tries to invoke
// the /ttt command. Its fields are populated straight from the data received after
// parsing the request form.
type RequestData struct {
	text        string // Represents the raw text input that follows /ttt
	channel     string // Channel ID
	userID      string // Current user's ID
	responseURL string
	username    string
	token       string //FIXME: not used?
	teamID      string //FIXME: not used?
	channelName string //FIXME: not used?
}

// GameHandler is the main handler that handles all incoming requests when a /ttt command is
// invoked. It parses the request into a RequestData object and sends the reqest off to be
// handled depending on what command it detects (start, move, display, or help).
func GameHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error has occurred: %v", err))
	}
	req := RequestData{
		text:        r.Form["text"],
		channel:     r.Form["channel_id"],
		userID:      r.Form["user_id"],
		responseURL: r.Form["response_url"],
		username:    r.Form["user_name"],
		token:       r.Form["token"],
		teamID:      r.Form["team_id"],
		channelName: r.Form["channel_name"],
	}

	inputList := strings.Split(req.text, " ")
	if len(inputList) <= 0 {
		fmt.Fprint(w, UsageError.Error())
		return
	}
	command := inputList[0]
	switch command {
	case "start":
		startGame(w, inputList, req)
	case "move":
		makeMove(w, inputList, req)
	case "display":
		handleDisplay(w, req)
	case "help":
		handleHelp(w, req)
	default:
		fmt.Fprint(w, UsageError.Error())
		return
	}
	return
}

func startGame(w http.ResponseWriter, inputList []string, req RequestData) {
	// TODO
}

func makeMove(w http.ResponseWriter, inputList []string, req RequestData) {
	// TODO
}

func handleDisplay(w http.ResponseWriter, inputList []string, req RequestData) {
	// TODO
}

func handleHelp(w http.ResponseWriter, inputList []string, req RequestData) {
	// TODO
}
