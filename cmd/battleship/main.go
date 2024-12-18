package main

import (
	"battleship/internal/TicTacToe"
)

func main() {
	// have a cli that asks user whether they want to be a server or a player
	// if they choose a player, let them join a game via ip
	// if they choose a server, ask them which game and load that game on the server
	tictactoe.CreateTicTacToeServer()
}
