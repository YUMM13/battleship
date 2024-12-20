package cmd

import (
	tictactoe "minigames/internal/TicTacToe"
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Generate a server.",
	Long: `Create a Tic Tac Toe, Sequence, or Battleship server.`,
	Run: createServer,
}

func init() {
	// add command to root command
	rootCmd.AddCommand(serverCmd)

	// add flags
	serverCmd.Flags().BoolP("tictactoe", "t", false, "Generate a Tic Tac Toe server.")
	serverCmd.Flags().BoolP("sequence", "s", false, "Generate a Sequence server.")
	serverCmd.Flags().BoolP("battleship", "b", false, "Generate a Battleship server.")
}

func createServer(cmd *cobra.Command, args []string) {
	// grab flags, ignore errors
	isTicTacToe, _ := cmd.Flags().GetBool("tictactoe")
	isSequence, _ := cmd.Flags().GetBool("sequence")
	isBattleship, _ := cmd.Flags().GetBool("battleship")

	if isTicTacToe && !isSequence && !isBattleship {
		tictactoe.CreateTicTacToeServer()
	} else if isSequence && !isTicTacToe && !isBattleship {
		fmt.Println("Sequence server is not yet implemented")
	} else if isBattleship && !isSequence && !isTicTacToe {
		fmt.Println("Battleship server is not yet implemented")
	} else {
		fmt.Println("use only one server flag!")
	}
}