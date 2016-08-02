package main

import (
	"errors"
	"fmt"
)

var HelpText = `Use /ttt to play a game of tic tac toe.
To start a game: /ttt start [@user]
To make a move: /ttt move [position]
To display current board: /ttt display
To cancel a current game: /ttt cancel`

// ----- ERRORS ----- //

var GenericError = errors.New(`Whoops! An error popped up out of nowhere. Try again, or try /ttt help.`)

var UsageError = errors.New(fmt.Sprintf("%s\nFor help: /ttt help.", HelpText))

var GameAlreadyExistsError = errors.New(`A game is already being played in this channel. Try another channel, or /ttt help for help.`)

var NoGameExistsError = errors.New(`No game is being played yet! Start one with /ttt start [@user], or try /ttt help for help.`)

var InvalidMoveError = errors.New(`That's not a valid move! Specify a position on the board using a row letter (A, B, C) and a column number (1, 2, 3).
	For example, to mark the bottom middle spot of the board: /ttt move C2`)

var PositionTakenError = errors.New(`That position is already taken!`)

var UserDoesntExistError = errors.New(`That user doesn't exist! Try again, or try /ttt help.`)

var InvalidTokenError = errors.New(`That's an invalid token. Which means you're an imposter, and you don't get to play!`)
