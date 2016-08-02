package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

func validToken(token string) bool {
	if token != authToken {
		return false
	}
	return true
}

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

type UsersResponse struct {
	Ok      bool     `json:"ok"`
	Members []Member `json:"members"`
	Error   string   `json:"error"`
}

type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getUsers(ctx context.Context) error {
	req, err := http.NewRequest("POST", "https://slack.com/api/users.list", nil)
	if err != nil {
		return fmt.Errorf("Couldn't create request to get users: %v", err)
	}
	values := req.URL.Query()
	values.Add("token", testToken)
	req.URL.RawQuery = values.Encode()

	c := urlfetch.Client(ctx)
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Couldn't do request to get users: %v", err)

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
	for _, m := range users.Members {
		Users[m.Name] = m.ID
	}
	return nil
}
