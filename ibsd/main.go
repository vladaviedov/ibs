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
	"time"
)

const version = "0.1.0"

type Report struct {
	Identifier string `json:"identifier"`
	MAC string `json:"mac"`
	IP string `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
	Passkey string `json:"passkey"`
}

type Settings struct {
	Identifier string `json:"identifier"`
	Server string `json:"server"`
	Passkey string `json:"passkey"`
	NetInterface string `json:"netInterface"`
}

func main() {
	// Command line args
	var configPath string
	if len(os.Args) == 2 {
		// Check options
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			printVersion()
			os.Exit(0)
		}
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			printUsage()
			os.Exit(0)
		}

		configPath = os.Args[1]
	} else {
		configPath = "/etc/ibsd/config.json"
	}

	// Read config
	var config Settings
	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open config file")
		fmt.Fprintln(os.Stderr, "Tried to access: " + configPath)
		os.Exit(1)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file")
		os.Exit(1)
	}

	for {
		sendReport(&config)
		time.Sleep(time.Minute * 5)
	}
}

func sendReport(config *Settings) {
	// Get net interface
	netInterface, err := net.InterfaceByName(config.NetInterface)
	if err != nil {
		log.Fatal("Net interface not available: ", err.Error())
	}

	// Make report
	report := Report{
		Identifier: config.Identifier,
		MAC: getMac(netInterface),
		IP: getIp(netInterface),
		Timestamp: time.Now(),
		Passkey: config.Passkey,
	}

	requestBody, err := json.Marshal(report)
	if err != nil {
		log.Fatal("Failed to serialize report: ", err.Error())
	}

	// Try https
	url := "https://" + config.Server
	res, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err == nil {
		if res.StatusCode != http.StatusOK {
			log.Println("Request returned error")
		}

		return
	}

	// Try http
	url = "http://" + config.Server
	res, err = http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Failed to reach server")
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("Request returned error")
	}
}

func getIp(netInterface *net.Interface) string {
	addrs, err := netInterface.Addrs()
	if err != nil {
		log.Fatal("Failed to get network address: ", err.Error())
	}

	for _, addr := range addrs {
		// Apparently this works
		if strings.Count(addr.String(), ":") < 2 {
			return strings.Split(addr.String(), "/")[0]
		}
	}

	log.Fatal("No IPv4 address found on this interface")
	return ""
}

func getMac(netInterface *net.Interface) string {
	// Verify MAC
	hwAddr := netInterface.HardwareAddr.String()
	mac, err := net.ParseMAC(hwAddr)
	if err != nil {
		log.Fatal("Failed to parse MAC address: ", err.Error())
	}

	return mac.String()
}

func printVersion() {
	fmt.Printf("IBS daemon v%s\n", version)
}

func printUsage() {
	fmt.Println("Usage: ibsd [opt] [configpath]")
	fmt.Println()
	fmt.Println("IBS daemon sends out device network information so that it can be found by IBS clients.") 
	fmt.Println("A report will be sent every 5 minutes when this program is running.")
	fmt.Println("By default /etc/ibsd/config.json will be parsed for settings.")
	fmt.Println("The used config file path can be changed by passing a filepath as an argument.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("%-15s: %s", "-v, --version", "Print command version\n")
	fmt.Printf("%-15s: %s", "-h, --help", "Print usage information\n")
}
