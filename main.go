package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
)

// Global state variables
var CurrentGames map[string]*Game    // Stores all the games in play currently, in all channels; maps channelID to Game
var ChannelUsers map[string][]string // Used for caching the users in each channel; maps channel ID to list of userIDs
var Users map[string]string          // List of users on this team; maps usernames to userIDs

func init() {
	CurrentGames = map[string]*Game{}
	ChannelUsers = map[string][]string{}
	Users = map[string]string{}
	http.HandleFunc("/play", GameHandler)
}

// RequestData represents an incoming request from a Slack channel when a user tries to invoke
// the /ttt command. Its fields are populated straight from the data received after
// parsing the request form. Only the fields used are part of the object here; any extraneous
// or unused fields are left out.
type RequestData struct {
	text        string // Represents the raw text input that follows /ttt
	channel     string // Channel ID
	userID      string // Current user's ID
	responseURL string // URL to send a response to Slack
	username    string // Current user's username
}

// ResponseData represents the data that should be sent back as a POST to the response URL
// provided in the request.
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
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, GenericError.Error())
		return
	}
	if !validToken(r.Form["token"][0]) {
		fmt.Fprint(w, InvalidTokenError.Error())
		return
	}

	ctx := appengine.NewContext(r)
	// Build a list of users on this team and in this channel if they're not already saved in memory.
	// NOTE: Unfortunately need to do this inside the handler because Google appengine requires
	// a current request context to make any http requests.
	err = getUserLists(ctx, r.Form["channel"][0])
	if err != nil {
		fmt.Fprint(w, GenericError.Error())
		return
	}
	// Parse request data. Assumes request is always well-formed, i.e. there is exactly one value
	// for each of the form keys
	requestData := RequestData{
		text:        r.Form["text"][0],
		channel:     r.Form["channel"][0],
		userID:      r.Form["user_id"][0],
		responseURL: r.Form["response_url"][0],
		username:    r.Form["user_name"][0],
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
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	// Generate JSON from the response and send it back
	json, err := response.getJSON()
	if err != nil {
		fmt.Fprintf(w, GenericError.Error())
		return
	}
	err = sendResponseData(requestData.responseURL, json, ctx)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	return
}
