package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Report struct {
	DeviceName string `json:"deviceName"`
	AddressReport string `json:"addressReport"`
	Passkey string `json:"passkey"`
}

type Status struct {
	Address string `json:"address"`
}

var router = mux.NewRouter().StrictSlash(true)
var deviceMap = make(map[string]string)
var passkey string;

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Invalid argument count. Usage: ibss <passkey>")
		return
	}
	passkey = os.Args[1];

	router.HandleFunc("/", home).Methods(http.MethodGet)
	router.HandleFunc("/report", report).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to home!")
}

func report(w http.ResponseWriter, r *http.Request) {
	// Parse report body
	var report Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Bad arguments")
		return
	}
	
	// Check passkey
	if report.Passkey != passkey {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Permission denied")
		return
	}

	// Create/update map
	_, ok := deviceMap[report.DeviceName]
	if !ok {
		router.HandleFunc(fmt.Sprintf("/%s", report.DeviceName), retrieveReport)
	}
	deviceMap[report.DeviceName] = report.AddressReport
}

func retrieveReport(w http.ResponseWriter, r *http.Request) {
	deviceName := r.URL.String()[1:]
	deviceAddr, ok := deviceMap[deviceName]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal error")
		return
	}

	status := Status{
		Address: deviceAddr,
	}
	json.NewEncoder(w).Encode(status)
}
