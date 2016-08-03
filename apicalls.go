package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// UsersResponse and User represent the response from a call to the Slack API to get a list of users
// for this team.
type UsersResponse struct {
	Ok    bool   `json:"ok"`
	Users []User `json:"members"`
	Error string `json:"error"`
}
type User struct {
	Name string
	ID   string
}

// ChannelUsersResponse and Channel represent the response from a call to the Slack API to get a list of users
// for the given channel.
type ChannelUsersResponse struct {
	Ok      bool    `json:"ok"`
	Channel Channel `json:"channel"`
	Error   string  `json:"error"`
}
type Channel struct {
	Members []string `json:"members"`
}

// makeAPICall is the generic method to make a call to the Slack API given a url to send the request to
// and a list of values to add to the URL. It uses the appengine urlfetch to create a client and do the
// request, which requires a context.Context from a currently executing http.Request.
func makeAPICall(url string, urlVals map[string]string, ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create request to get users: %v", err)
	}
	values := req.URL.Query()
	for key, val := range urlVals {
		values.Add(key, val)
	}
	req.URL.RawQuery = values.Encode()

	c := urlfetch.Client(ctx)
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't do request to get users: %v", err)
	}
	return resp, nil
}

// getUsers makes an API call to get a list of users for this team, then parses the response and stores
// the list of users in the global Users map defined in main.go.
func getUsers(ctx context.Context) error {
	resp, err := makeAPICall("https://slack.com/api/users.list", map[string]string{"token": testToken}, ctx)
	if err != nil {
		return fmt.Errorf("Failed to make API call to get users: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Couldn't read body of response: %v", err)
	}
	var users UsersResponse
	err = json.Unmarshal(body, &users)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshal response: %v", err)
	}
	if !users.Ok {
		return fmt.Errorf("Error getting users: %s", users.Error)
	}
	for _, m := range users.Users {
		Users[m.Name] = m.ID
	}
	return nil
}

// getChannelUsers makes an API call to get a list of users for the given channel, then parses the response and stores
// the list of userIDs in the global ChannelUsers map defined in main.go.
func getChannelUsers(channelID string, ctx context.Context) error {
	vals := map[string]string{
		"token":   testToken,
		"channel": channelID,
	}
	resp, err := makeAPICall("https://slack.com/api/channels.info", vals, ctx)
	if err != nil {
		return fmt.Errorf("Failed to make API call to get channel users: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Couldn't read body of response: %v", err)
	}
	var users ChannelUsersResponse
	err = json.Unmarshal(body, &users)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshal response: %v", err)
	}
	if !users.Ok {
		return fmt.Errorf("Error getting channel users: %s", users.Error)
	}
	ChannelUsers[channelID] = users.Channel.Members
	return nil
}
