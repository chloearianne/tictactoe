# tictactoe

A Slack custom integration for playing a game of tic tac toe within a channel.

## Usage
- `/ttt start [@user]` challenges @user to a game. Anyone can run this command.
- `/ttt display` displays the current status of the board and whose turn it currently is. Anyone can run this command.
- `/ttt move [space on board]` is used to make a move, marking the corresponding space with an X or an O depending on the user. Only the current user can run this command.
- `/ttt cancel` is used to cancel a game. Only the players of the game can run this command.
- `/ttt help` displays information about how to play the game. Anyone can run this command.

## Gameplay details
Positions on the board are represented by two characters: the first, a letter indicating the row (A, B, or C), and the second, a number indicating the column (1, 2, or 3). For example "/ttt move C2" would mark the bottom row, middle spot.
For rules of tic tac toe, see https://en.wikipedia.org/wiki/Tic-tac-toe.

Only one game is allowed be active per channel.
