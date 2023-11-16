package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Status struct {
	Address string `json:"address"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Invalid argument count. Usage: ibss <passkey>")
		return
	}

	// Make request
	url := fmt.Sprintf("http://localhost:8080/%s", os.Args[1])
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln("Failed to get response from server.")
		return
	}

	// Check status code
	if response.StatusCode != http.StatusOK {
		log.Fatalln("Device not found.")
		return
	}

	// Get device ip
	var status Status
	err = json.NewDecoder(response.Body).Decode(&status)
	if err != nil {
		log.Fatalln("Failed to parse response.")
		return
	}

	fmt.Println(status.Address)
}
