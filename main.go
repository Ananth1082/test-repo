package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed index.html
var html string

func main() {
	fmt.Println("Serving HTML content: ", html)

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recieved request from client")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, html); err != nil {
			fmt.Println("Error writing response: ", err)
		}
	})
	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8800", router); err != nil {
		panic("Server couldnt start")
	}
}
