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
	if len(os.Args) != 3 {
		log.Fatalln("Invalid argument count. Usage: ibsc <server> <device name>")
		return
	}

	// Make request
	url := fmt.Sprintf("http://%s/%s", os.Args[1], os.Args[2])
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
