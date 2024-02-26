package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const pattern = `\S*\.ibs\b`

type Settings struct {
	Server string `json:"server"`
}

func main() {
	if len(os.Args) == 1 {
		printStatus()
		os.Exit(0)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "No home directory :(")
		os.Exit(1)
	}

	// Read config
	configPath := homeDir + "/.ibsc_conf"
	var config Settings
	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open config file")
		fmt.Fprintln(os.Stderr, "Tried to access: ", configPath)
		os.Exit(1)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file")
		os.Exit(1)
	}

	reg, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad matching pattern")
		os.Exit(1)
	}

	// Find ibs addresses
	for i, arg := range os.Args {
		// Skip ibsc
		if i == 0 {
			continue;
		}

		for {
			match := reg.FindString(arg)
			if match == "" {
				break
			}

			ip := resolveDomain(&config, match)
			arg = strings.ReplaceAll(arg, match, ip)
		}

		os.Args[i] = arg
	}

	fmt.Println(strings.Join(os.Args[1:], " "))
}

func printStatus() {
	
}

func resolveDomain(config *Settings, domain string) string {
	identifier := strings.TrimSuffix(domain, ".ibs")

	// Try https
	url := "https://" + config.Server + "/dns/" + identifier
	res, err := http.Get(url)
	if err != nil {
		// Try http
		url = "http://" + config.Server + "/dns/" + identifier
		res, err = http.Get(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to reach server")
			os.Exit(1)
		}
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Server failed to resolve ", domain)
		os.Exit(1)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get response body: ", err.Error())
		os.Exit(1)
	}

	return string(bytes)
}
