package tictactoe

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aquasecurity/table"
)

// type for game state
type GameState struct {
	Board        [3][3]string
	CurrentTurn  string
	Winner       string
	PlayerTokens map[string]string // maps players to session tokens to ensure players can only make play for themselves
	PlayerNames  map[string]string // maps custom player names to their role in the game (x or O)
}

// set initial board state
var state = GameState{
	Board: [3][3]string{{"0", "1", "2"},
						{"3", "4", "5"},
						{"6", "7", "8"}},
	CurrentTurn:  "X",
	Winner:       "", // no winner
	PlayerTokens: make(map[string]string),
	PlayerNames:  make(map[string]string),
}

// create tic tac toe server
func CreateTicTacToeServer() {
	// create a new mux server
	mux := http.NewServeMux()

	// create the paths that will be handled by the server
	mux.HandleFunc("GET /join/", joinHandler)
	mux.HandleFunc("GET /move/", moveHandler)
	mux.HandleFunc("GET /reset/", resetHandler)

	// start the server and check for errors
	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		fmt.Println("Error starting Tic Tac Toe server:", err)
	}
}

// resets the board to an initial value, but keeps the original players
func setBoardState() {
	state.Board = [3][3]string{{"0", "1", "2"},
		{"3", "4", "5"},
		{"6", "7", "8"}}
	state.CurrentTurn = "X"
	state.Winner = ""
}

// gets the board's current state and sends it to a buffer. The returned buffer will be used in various
// HTTP responses
func getBoardState() bytes.Buffer {
	// create table
	var buffer bytes.Buffer
	table := table.New(&buffer)
	table.AddHeaders("TIC", "TAC", "TOE")

	for _, p := range state.Board {
		table.AddRow(p[0], p[1], p[2])
	}

	// render the table to the buffer
	table.Render()

	return buffer
}

// handle join requests
func joinHandler(w http.ResponseWriter, r *http.Request) {
	// get player from request and do error checking on it
	player := r.URL.Query().Get("player")
	if player == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}
	if _, exists := state.PlayerNames[player]; exists {
		http.Error(w, "Player already joined", http.StatusConflict)
		return
	}

	// temp until actual security is implemented
	// will prob have server generate token and the client will hold it and send it automatically
	state.PlayerTokens[player] = "123"

	// assign roles
	if len(state.PlayerNames) == 0 {
		state.PlayerNames[player] = "X"
	} else if len(state.PlayerNames) == 1 {
		state.PlayerNames[player] = "O"
	} else {
		http.Error(w, "Lobby is full", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Server received join request, welcome %s", player)
}

// handle move requests
func moveHandler(w http.ResponseWriter, r *http.Request) {
	// check to make sure the correct player sent it
	player := r.Header.Get("Player")
	token := r.Header.Get("Authorization")

	err := validatePlayer(player, token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// check to make sure it is the players turn
	if state.PlayerNames[player] != state.CurrentTurn {
		http.Error(w, "It is not your turn", http.StatusUnauthorized)
		return
	}
	
	// gets the position in the move request
	pos, err := strconv.Atoi(r.URL.Query().Get("position"))

	// error checking
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %q is not an integer.", r.PathValue("position")), http.StatusBadRequest)
		return
	}

	if pos < 0 || pos > 8 {
		http.Error(w, fmt.Sprintf("Error: %q is not a valid position.", pos), http.StatusBadRequest)
		return
	}

	// place the piece
	if state.CurrentTurn == "X" {
		state.Board[pos/3][pos%3] = "X"
		state.CurrentTurn = "O"
	} else {
		state.Board[pos/3][pos%3] = "O"
		state.CurrentTurn = "X"
	}

	// get the board state after move is made
	buf := getBoardState()

	// Send the buffer's content as the HTTP response
	w.Header().Set("Content-Type", "text/plain")
	if checkForWinner() {
		fmt.Fprintf(w, "3 in a row! %s wins. Sending board state:\n\n%s", player, buf.String())
	} else {
		fmt.Fprintf(w, "Server received move request. Sending board state:\n\n%s", buf.String())
	}
}

// helper function that validates the player
func validatePlayer(player string, auth string) error {
	if state.PlayerTokens[player] != auth {
        return fmt.Errorf("invalid token")
	}
	return nil
}

// helper function that checks if there is a winner
func checkForWinner() bool {
	if state.Board[0][0] == state.Board[0][1] && state.Board[0][1] == state.Board[0][2] {
		return true
	} else if state.Board[1][0] == state.Board[1][1] && state.Board[1][1] == state.Board[1][2] {
		return true
	} else if state.Board[2][0] == state.Board[2][1] && state.Board[2][1] == state.Board[2][2] {
		return true
	} else if state.Board[0][0] == state.Board[1][0] && state.Board[1][0] == state.Board[2][0] {
		return true
	} else if state.Board[0][1] == state.Board[1][1] && state.Board[1][1] == state.Board[2][1] {
		return true
	} else if state.Board[0][2] == state.Board[1][2] && state.Board[1][2] == state.Board[2][2] {
		return true
	} else if state.Board[0][0] == state.Board[1][1] && state.Board[1][1] == state.Board[2][2] {
		return true
	} else if state.Board[0][2] == state.Board[1][1] && state.Board[1][1] == state.Board[2][0] {
		return true
	} else {
		return false
	}
}

// handle reset requests
func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server received reset request. Resetting game.")
	setBoardState()
}
