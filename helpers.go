package main

import (
	"bytes"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

// getUserLists is used to build lists of users on this team and in this channel if they're
// not already saved in memory.
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

func moveIsValid(move string) bool {
	for _, pos := range boardPositions {
		if strings.ToUpper(move) == pos {
			return true
		}
	}
	return false
}
