package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Report struct {
	DeviceName string `json:"deviceName"`
	AddressReport string `json:"addressReport"`
	Passkey string `json:"passkey"`
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("Invalid argument count. Usage: ibsd <server> <device name> <passkey>")
		return
	}

	report := Report{
		DeviceName: os.Args[2],
		AddressReport: getIp(),
		Passkey: os.Args[3],
	}
	url := fmt.Sprintf("http://%s/report", os.Args[1])

	// Send request
	jsonReport, jsonErr := json.Marshal(report)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReport))
	if err != nil {
		log.Fatal(err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("Request returned error.")
		return
	}
}

func getIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fullAddr := conn.LocalAddr().String()
	return strings.Split(fullAddr, ":")[0]
}
