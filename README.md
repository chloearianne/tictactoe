# tictactoe

A Slack custom integration for playing a game of tic tac toe within a channel.

## Usage
- `/ttt start [@user]` challenges @user to a game.
- `/ttt display` displays the current status of the board and whose turn it currently is.
- `/ttt move [space on board]` is used to make a move, marking the corresponding space with an X or an O depending on the user.

## Gameplay details
The board is represented in a grid fashion, with rows labeled A, B, and C, and columns labeled 1, 2, and 3. For example, to mark the bottom row, center spot on your turn, you would use the command `/ttt move C2`. 
Only one game is allowed be active per channel.
