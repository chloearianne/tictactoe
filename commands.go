package main

import (
	"fmt"
	"strings"
)

// handleStart provides the functionality for starting a game. It makes sure the command and
// other arguments (given in inputList) are formatted correctly, that the provided user exists and
// is in the current channel, and then creates and saves a new Game with the players' info. It posts
// a public response in the channel saying that a user has been challenged to a game.
func handleStart(inputList []string, req RequestData) (*ResponseData, error) {
	// Ensure correct number of arguments, ex. start @bob
	if len(inputList) != 2 {
		return nil, UsageError
	}
	// Check if a game is already being played in this channel
	if _, exists := CurrentGames[req.channel]; exists {
		return nil, GameAlreadyExistsError
	}
	// Look up user to be challenged, make sure they exist and are in this channel
	challengedUser := inputList[1]
	if !strings.HasPrefix(challengedUser, "@") {
		return nil, UsageError
	}
	challengedUser = strings.TrimPrefix(challengedUser, "@")
	challengedUserID, ok := Users[challengedUser]
	if !ok {
		return nil, UserDoesntExistError
	}
	userList, ok := ChannelUsers[req.channel]
	if ok {
		found := false
		for _, user := range userList {
			if user == challengedUserID {
				found = true
			}
		}
		if !found {
			return nil, UserNotInChannelError
		}
	}
	// Create players and new game
	p2 := Player{
		Name: challengedUser,
		ID:   challengedUserID,
		Mark: X,
	}
	p1 := Player{
		Name: req.username,
		ID:   req.userID,
		Mark: O,
	}
	CurrentGames[req.channel] = New(p1, p2)

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("<@%s|%s>, %s has challenged you to a game! Your move. \n%s", challengedUserID, challengedUser, req.username, emptyBoardDisplay),
	}

	return &response, nil
}

// handleMove provides the functionality for users to make moves in the game. It checks for proper
// formatting of the command and its arguments passed in through inputList, verifies that a game is
// in progress and that it's the current user's turn, verifies that the move is valid, and then makes
// the move. It then checks if the game has been won or tied, and if not, flips who the current player
// is and sends a response to display the updated board.
func handleMove(inputList []string, req RequestData) (*ResponseData, error) {
	// Ensure correct number of arguments, ex. move A1
	if len(inputList) != 2 {
		return nil, UsageError
	}
	// Check if a game is being played in this channel
	game, ok := CurrentGames[req.channel]
	if !ok {
		return nil, NoGameExistsError
	}
	// Check if it's the current user's turn
	if req.userID != game.CurrentPlayer.ID {
		if req.userID != game.Player1.ID && req.userID != game.Player2.ID {
			return nil, NotAuthorizedError
		}
		return nil, NotYourTurnError
	}
	// Check validity of move
	move := strings.ToUpper(inputList[1])
	if !moveIsValid(move) {
		return nil, InvalidMoveError
	}
	if game.Board[move] != empty {
		return nil, PositionTakenError
	}
	game.Board[move] = game.CurrentPlayer.Mark
	// Deal with end of game scenarios
	if game.HasWinner() {
		text := fmt.Sprintf("%s\n%s has won the game!\nGame over. (If the result displeases you, you could always play another game...)", game.Display(), req.username)
		response := ResponseData{
			ResponseType: "in_channel",
			Text:         text,
		}
		delete(CurrentGames, req.channel)
		return &response, nil
	}
	if game.IsOver() {
		text := fmt.Sprintf("%s\nYou've tied! Your skills are matched, apparently.\n (You could play another game to find out for sure...)", game.Display())
		response := ResponseData{
			ResponseType: "in_channel",
			Text:         text,
		}
		delete(CurrentGames, req.channel)
		return &response, nil
	}
	// Switch who the current player is
	if game.CurrentPlayer == game.Player1 {
		game.CurrentPlayer = game.Player2
	} else {
		game.CurrentPlayer = game.Player1
	}
	CurrentGames[req.channel] = game

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s\nIt's %s's turn to make a move.", game.Display(), game.CurrentPlayer.Name),
	}
	return &response, nil
}

// handleDisplay provides the functionality for any user in the channel to display the
// current board and whose turn it is.
func handleDisplay(inputList []string, req RequestData) (*ResponseData, error) {
	// Ensure correct number of arguments (just display)
	if len(inputList) != 1 {
		return nil, UsageError
	}
	// Check if a game is being played in this channel
	game, ok := CurrentGames[req.channel]
	if !ok {
		return nil, NoGameExistsError
	}
	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s\nIt's %s's turn to make a move.", game.Display(), game.CurrentPlayer.Name),
	}
	return &response, nil
}

// handleHelp returns the help text defined in constants.go.
func handleHelp() (*ResponseData, error) {
	response := ResponseData{
		ResponseType: "ephemeral",
		Text:         HelpText,
	}
	return &response, nil
}

// handleCancel provides functionality for cancelling a gmae that is currently in session.
// It checks if a game exists in the channel and if the user is authorized to cancel it
// (i.e. the user is one of the two players) and then removes it from the global
// list of games.
func handleCancel(inputList []string, req RequestData) (*ResponseData, error) {
	// Ensure correct number of arguments (just cancel)
	if len(inputList) != 1 {
		return nil, UsageError
	}
	// Check if a game is being played in this channel
	game, ok := CurrentGames[req.channel]
	if !ok {
		return nil, NoGameExistsError
	}
	// Check if the user trying to cancel the game is one of the players
	if req.userID != game.Player1.ID && req.userID != game.Player2.ID {
		return nil, NotAuthorizedError
	}
	delete(CurrentGames, req.channel)

	response := ResponseData{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s has cancelled the current game. Perhaps a rematch later.", req.username),
	}
	return &response, nil
}
