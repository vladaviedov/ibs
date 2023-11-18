package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Invalid argument count. Usage: ibss <passkey>")
		return
	}
	Passkey = os.Args[1]

	Router.HandleFunc("/", home).Methods(http.MethodGet)
	Router.HandleFunc("/report", ProcessReport).Methods(http.MethodPost)
	Router.HandleFunc("/list", ListDevices).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to home!")
}
