package main

import (
	"fmt"
	"strings"
)

func handleStart(inputList []string, req RequestData) (*ResponseData, error) {
	if len(inputList) != 2 {
		return nil, UsageError
	}
	p1 := Player{
		Name: req.username,
		ID:   req.userID,
	}
	// Extract name of user to be challenged, and look up their unique user ID
	challengedUser := inputList[1]
	challengedUser = strings.TrimPrefix(challengedUser, "@")
	challengedUserID, ok := Users[challengedUser]
	if !ok {
		return nil, UserDoesntExistError
	}
	p2 := Player{
		Name: challengedUser,
		ID:   challengedUserID,
	}

	// Check if a game is already being played in this channel
	if _, exists := CurrentGames[req.channel]; exists {
		return nil, GameAlreadyExistsError
	}

	CurrentGames[req.channel] = *New(p1, p2)

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("<@%s|%s>, %s has challenged you to a game! To accept this noble challenge, make the first move.", challengedUserID, challengedUser, req.username),
	}

	return &response, nil
}

func handleMove(inputList []string, req RequestData) (*ResponseData, error) {
	//TODO
	return nil, nil
}

func handleDisplay(inputList []string, req RequestData) (*ResponseData, error) {
	if len(inputList) != 1 {
		return nil, UsageError
	}
	game, ok := CurrentGames[req.channel]
	if !ok {
		return nil, NoGameExistsError
	}

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         game.Display(),
	}
	return &response, nil
}

func handleHelp() (*ResponseData, error) {
	response := ResponseData{
		ResponseType: "ephemeral",
		Text:         HelpText,
	}
	return &response, nil
}

func handleCancel(inputList []string, req RequestData) (*ResponseData, error) {
	if len(inputList) != 1 {
		return nil, UsageError
	}

	game, ok := CurrentGames[req.channel]
	if !ok {
		return nil, NoGameExistsError
	}
	if req.userID != game.Player1.ID && req.userID != game.Player2.ID {
		return nil, NotAuthorizedError
	}
	delete(CurrentGames, req.channel)

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s has cancelled the current game. What a shame.", req.username),
	}
	return &response, nil
}
