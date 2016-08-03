package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// Global state variables
var CurrentGames map[string]*Game    // Stores all the games in play currently (for all channels); maps channelID to Game
var ChannelUsers map[string][]string // Used for caching the users in each channel; maps channelID to list of userIDs
var Users map[string]string          // List of users on this Slack team; maps usernames to userIDs

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
	text        string // Raw text input that follows /ttt
	channel     string // Channel ID
	userID      string // Current user's ID
	responseURL string // URL to send a response to Slack
	username    string // Current user's username
}

// ResponseData represents the data that should be sent back as a POST to the response URL
// provided in the request data.
type ResponseData struct {
	ResponseType string `json:"response_type"` // values: in_channel, ephemeral
	Text         string `json:"text"`          // text to be written in the message
}

func (r ResponseData) getJSON() ([]byte, error) {
	jsonResponse, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// GameHandler is the main handler that processes all incoming requests from when a /ttt command is
// invoked. It parses the request into a RequestData object and sends the request off to be
// handled depending on what command it detects (start, move, display, or help).
func GameHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, GenericError.Error())
		return
	}

	if !validFormValues(r.Form) {
		fmt.Fprint(w, GenericError.Error())
		return
	}
	if !validToken(r.Form["token"][0]) {
		fmt.Fprint(w, InvalidTokenError.Error())
		return
	}

	ctx := appengine.NewContext(r)
	// Build a list of users on this team and in this channel if they're not already saved in memory.
	// NOTE: Need to do this inside the handler because Google appengine requires
	// a current appengine http.Request context to create and execute http requests.
	err = getUserLists(ctx, r.Form["channel_id"][0])
	if err != nil {
		fmt.Fprint(w, GenericError.Error())
		return
	}

	requestData := RequestData{
		text:        r.Form["text"][0],
		channel:     r.Form["channel_id"][0],
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
		fmt.Fprint(w, err.Error())
		return
	}

	json, err := response.getJSON()
	if err != nil {
		fmt.Fprint(w, GenericError.Error())
		return
	}
	err = sendResponseData(requestData.responseURL, json, ctx)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	return
}

// validFormValues verifies if the form passed in (from an http request)
// have the proper values we need to process the request
func validFormValues(form url.Values) bool {
	if len(form["text"]) != 1 ||
		len(form["token"]) != 1 ||
		len(form["channel_id"]) != 1 ||
		len(form["user_id"]) != 1 ||
		len(form["response_url"]) != 1 ||
		len(form["user_name"]) != 1 {
		return false
	}
	return true
}

// getUserLists is used to build lists of users on this team and in this channel if they're
// not already saved in memory. To get this data, it uses API calls (see getUsers and getChannelUsers
// in  apicalls.go).
func getUserLists(ctx context.Context, channelID string) error {
	if len(Users) == 0 {
		err := getUsers(ctx)
		if err != nil {
			return err
		}
	}
	_, ok := ChannelUsers[channelID]
	if !ok {
		err := getChannelUsers(channelID, ctx)
		if err != nil {
			return err
		}
	}
	return nil

}

func validToken(token string) bool {
	if token != authToken {
		return false
	}
	return true
}

// sendResponseData sends JSON responses back to Slack after the /ttt command
// has been handled.
func sendResponseData(url string, json []byte, ctx context.Context) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return GenericError
	}
	req.Header.Add("Content-Type", "application/json")
	c := urlfetch.Client(ctx)
	_, err = c.Do(req)
	if err != nil {
		return GenericError
	}
	return nil
}
