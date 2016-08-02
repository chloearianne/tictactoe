package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
)

var CurrentGames map[string]Game
var Users map[string]string

func init() {
	CurrentGames = map[string]Game{}
	Users = map[string]string{}
	http.HandleFunc("/play", GameHandler)
}

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

// ResponseData represents the data that should be sent back as a POST to the response URL
// provided in the request. It
type ResponseData struct {
	ResponseType string `json:"response_type"` // values: in_channel, ephemeral
	Text         string `json:"text"`
}

func (r ResponseData) getJSON() ([]byte, error) {
	jsonResponse, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// GameHandler is the main handler that handles all incoming requests when a /ttt command is
// invoked. It parses the request into a RequestData object and sends the reqest off to be
// handled depending on what command it detects (start, move, display, or help).
func GameHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// Build a list of users if it's not already saved in memory.
	// Unfortunately need to do this inside the handler because Google appengine requires
	// a current request context to make any http requests
	// TOOD: look into better ways of doing this
	if len(Users) == 0 {
		err := getUsers(ctx)
		if err != nil {
			fmt.Fprintf(w, GenericError.Error())
			return
		}
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, GenericError.Error())
		return
	}

	if !validToken(r.Form["token"][0]) {
		fmt.Fprint(w, InvalidTokenError.Error())
		return
	}
	//TODO: parse request data using extenral package
	// Assumes request is always well-formed
	requestData := RequestData{
		text:        r.Form["text"][0],
		channel:     r.Form["channel_id"][0],
		userID:      r.Form["user_id"][0],
		responseURL: r.Form["response_url"][0],
		username:    r.Form["user_name"][0],
		token:       r.Form["token"][0],
		teamID:      r.Form["team_id"][0],
		channelName: r.Form["channel_name"][0],
	}

	inputList := strings.Split(requestData.text, " ")
	if requestData.text == "" || len(inputList) <= 0 {
		fmt.Fprint(w, UsageError.Error())
		return
	}
	var response *ResponseData
	command := inputList[0]
	switch command {
	case "start":
		response, err = handleStart(inputList, requestData)
	case "move":
		response, err = handleMove(inputList, requestData)
	case "display":
		response, err = handleDisplay(inputList, requestData)
	case "cancel":
		response, err = handleCancel(inputList, requestData)
	case "help":
		response, err = handleHelp()
	default:
		fmt.Fprint(w, UsageError.Error())
		return
	}

	// If handling the command produces an error, just return that
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	// Otherwise, generate JSON from the response
	json, err := response.getJSON()
	if err != nil {
		fmt.Fprintf(w, GenericError.Error())
		return
	}

	// Send the generated JSON back in a response
	err = sendResponseData(requestData.responseURL, json, ctx)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	return
}
