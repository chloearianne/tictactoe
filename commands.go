package main

import (
	"fmt"
	"strings"
)

func handleStart(inputList []string, req RequestData) (*ResponseData, error) {
	if len(inputList) != 2 {
		return nil, UsageError
	}

	// Extract name of user to be challenged, and look up their unique user ID
	challengedUser := inputList[1]
	challengedUser = strings.TrimPrefix(challengedUser, "@")
	challengedUserID, ok := Users[challengedUser]
	if !ok {
		return nil, UserDoesntExistError
	}

	// Check if a game is already being played in this channel
	if _, exists := CurrentGames[req.channel]; exists {
		return nil, GameAlreadyExistsError
	}
	game := New([]string{req.username, challengedUser}, []string{req.userID, challengedUserID}, challengedUserID)
	CurrentGames[req.channel] = *game

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("<@%s|%s>, %s has challenged you to a duel! To accept this noble challenge, make the first move.", challengedUserID, challengedUser, req.username),
	}

	return &response, nil
}

func handleMove(inputList []string, req RequestData) (*ResponseData, error) {
	//TODO
	return nil, nil
}

func handleDisplay(inputList []string, req RequestData) (*ResponseData, error) {
	// TODO
	return nil, nil
}

func handleHelp() (*ResponseData, error) {
	resp := ResponseData{
		ResponseType: "ephemeral",
		Text:         HelpText,
	}
	return &resp, nil
}

func handleCancel(inputList []string, req RequestData) (*ResponseData, error) {
	if _, ok := CurrentGames[req.channel]; !ok {
		return nil, NoGameExistsError
	}
	delete(CurrentGames, req.channel)

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s has cancelled the current game. What a shame.", req.username),
	}
	return &response, nil
}
