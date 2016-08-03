package main

import (
	"errors"
	"fmt"
)

// ------ GAME ------ //

// Valid Positions
var A1 = "A1"
var B1 = "B1"
var C1 = "C1"
var A2 = "A2"
var B2 = "B2"
var C2 = "C2"
var A3 = "A3"
var B3 = "B3"
var C3 = "C3"

var X = "X"
var O = "0" // Using 0 instead of O because the 0 is skinnier and makes the board layout look better
var empty = "..."

var boardPositions = []string{A1, A2, A3, B1, B2, B3, C1, C2, C3}

var emptyBoard = map[string]string{A1: empty, B1: empty, C1: empty, A2: empty, B2: empty, C2: empty, A3: empty, B3: empty, C3: empty}

var emptyBoardDisplay = `      1    2    3
A  ... | ... | ...
B  ... | ... | ...
C  ... | ... | ...`

// ------ USEFUL MESSAGES ------ //

var Usage = `To start a game: /ttt start [@user]
To make a move: /ttt move [position]
To display current board: /ttt display
To cancel a current game: /ttt cancel`

var HelpText = fmt.Sprintf(`Use /ttt to play a game of tic tac toe.
%s
Positions on the board are represented by two characters: the first, a letter indicating the row (A, B, or C), and the second, a number indicating the column (1, 2, or 3). For example "/ttt move C2" would result in this mark:

      1    2    3
A  ... | ... | ...
B  ... | ... | ...
C  ... | X | ...

For rules of tic tac toe, see <https://en.wikipedia.org/wiki/Tic-tac-toe>.`, Usage)

// ----- ERRORS ----- //

var GenericError = errors.New(`Whoops! An error popped up out of nowhere. Try again, or try /ttt help.`)
var UsageError = errors.New(fmt.Sprintf("Hmm, I didn't quite understand that.\n%s\nFor help: /ttt help.", Usage))
var GameAlreadyExistsError = errors.New(`A game is already being played in this channel. Try another channel, or /ttt help for help.`)
var NoGameExistsError = errors.New(`No game is being played yet! Start one with /ttt start [@user], or try /ttt help for help.`)
var InvalidMoveError = errors.New(`That's not a valid move! Specify a position on the board using a row letter (A, B, C) and a column number (1, 2, 3). For example, to mark the bottom middle spot of the board: /ttt move C2. Try /ttt help for more information.`)
var PositionTakenError = errors.New(`That position is already taken. It was a good one, I know, but try another.`)
var UserDoesntExistError = errors.New(`That user doesn't exist on this team! Try again, or try /ttt help.`)
var UserNotInChannelError = errors.New(`That user isn't in this channel. Invite them to join first, or play in another channel.`)
var InvalidTokenError = errors.New(`That's an invalid token. Which means you're an imposter, and you don't get to play!`)
var NotAuthorizedError = errors.New(`You can't do that, you're not even playing in this game!`)
var NotYourTurnError = errors.New(`It's not your turn.`)
