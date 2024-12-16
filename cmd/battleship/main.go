package main

import (
	"fmt"
	"net/http"
)

func main() {
	// create a new mux server
	mux := http.NewServeMux()

	// create the paths that will be handled by the server
	mux.HandleFunc("GET /path/", homeHandler)

	fmt.Println("Server is listening...")

	// start the server and check for errors
	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// when someone sends a GET request to /path/, this function is called
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}