package main

import (
	"net/http"
	"os"
	"tutorials/rotor-controller/rotor"
	"tutorials/rotor-controller/schedule"
)

func main() {
	http.HandleFunc("/rotor", rotor.HandleFunc)
	http.HandleFunc("/schedule", schedule.PassesHandleFunc)
	http.HandleFunc("/schedule/", schedule.PassHandleFunc)
	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
