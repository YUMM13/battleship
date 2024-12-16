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
	Board       [3][3]string
	CurrentTurn string
	Winner      string
	PlayerList  map[string]bool
}

// set initial board state
var state = GameState{
				Board:       [3][3]string{{"0", "1", "2"},
										  {"3", "4", "5"},
										  {"6", "7", "8"}},
				CurrentTurn: "X",
				Winner:      "", // no winner
				PlayerList:  make(map[string]bool),
			}

// create tic tac toe server
func CreateTicTacToeServer() {
	// create a new mux server
	mux := http.NewServeMux()

	// create the paths that will be handled by the server
	mux.HandleFunc("GET /join/{player}", joinHandler)
	mux.HandleFunc("GET /move/{position}", moveHandler)
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
	fmt.Fprintf(w, "Server received join request")
}

// handle move requests
func moveHandler(w http.ResponseWriter, r *http.Request) {
	// gets the position in the move request
	pos, err := strconv.Atoi(r.PathValue("position"))

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
		state.Board[pos / 3][pos % 3] = "X"
		state.CurrentTurn = "O"
	} else {
		state.Board[pos / 3][pos % 3] = "O"
		state.CurrentTurn = "X"
	}

	// get the board state after move is made
	buf := getBoardState()

	// Send the buffer's content as the HTTP response
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Server received move request. Sending board state:\n\n%s", buf.String())
}


// handle reset requests
func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server received reset request. Resetting game.")
	setBoardState()
}
