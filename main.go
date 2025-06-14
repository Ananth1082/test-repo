package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed index.html
var html string

func main() {
	fmt.Println("Serving HTML content: ", html)

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recieved request from client")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, html); err != nil {
			log.Println("Error writing response: ", err)
		}
	})

	router.HandleFunc("GET /message", func(w http.ResponseWriter, r *http.Request) {
		if fr, err := os.ReadFile("messages.log"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "<p>No messages!</p>")
		} else {
			messsages := string(fr)
			str := `<table border="1">`
			for i, msg := range strings.Split(messsages, "\n") {
				if msg == "" {
					continue
				}
				str += fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, msg)
			}
			str += "</table>"
			fmt.Fprintln(w, str)
		}
	})

	router.HandleFunc("POST /message", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body: ", err)
			return
		}

		fileWriter, err := os.OpenFile("messages.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening log file: ", err)
			return
		}
		if _, err := fmt.Fprintln(fileWriter, string(data)); err != nil {
			log.Println("Error writing logs: ", err)
		}
	})

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8134", router); err != nil {
		panic("Server couldnt start")
	}
}
