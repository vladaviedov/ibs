package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

const pattern = `[a-zA-Z0-9][a-zA-Z0-9-]*\.((ibs)|(IBS))\b`
const version = "0.2.0"

type Settings struct {
	Server string `json:"server"`
	WithShell bool `json:"withShell"`
	ShowCommand bool `json:"showCommand"`
}

func main() {
	if len(os.Args) == 1 {
		printStatus()
		os.Exit(0)
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		printVersion()
		os.Exit(0)
	}
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		os.Exit(0)
	}

	reg, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad matching pattern")
		os.Exit(1)
	}

	config := loadConfig()

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

	var bin string
	var argv []string
	if !config.WithShell {
		// Find needed binary, since we don't have execvp
		bin, err = exec.LookPath(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		argv = os.Args[1:]
	} else {
		// Try to get from env or fallback to sh
		bin = os.Getenv("SHELL")
		if bin == "" {
			fmt.Fprintln(os.Stderr, "Warning: $SHELL is not set, trying to find sh...")
			bin, err = exec.LookPath("sh")
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
		argv = append([]string{bin, "-c"}, strings.Join(os.Args[1:], " "))
	}

	// Run exec
	if config.ShowCommand {
		fmt.Println(strings.Join(argv, " "))
	}
	err = syscall.Exec(bin, argv, os.Environ())

	// Exec failed
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func loadConfig() Settings {
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
		fmt.Fprintln(os.Stderr, "Tried to access:", configPath)
		os.Exit(1)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file")
		os.Exit(1)
	}

	return config
}

func printStatus() {
	config := loadConfig()
	fmt.Println("Server:", config.Server)
	fmt.Println("--------------------")

	// Ping HTTPS
	url := "https://" + config.Server + "/ping"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("https: not reachable")
	} else if res.StatusCode != http.StatusOK {
		fmt.Println("https: server error")
	} else {
		fmt.Println("https: good")
	}

	// Ping HTTP
	url = "http://" + config.Server + "/ping"
	res, err = http.Get(url)
	if err != nil {
		fmt.Println("http:  not reachable")
	} else if res.StatusCode != http.StatusOK {
		fmt.Println("http:  server error")
	} else {
		fmt.Println("http:  good")
	}
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
		fmt.Fprintln(os.Stderr, "Server failed to resolve", domain)
		os.Exit(1)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get response body:", err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func printVersion() {
	fmt.Printf("IBS client v%s\n", version)
}

func printUsage() {
	fmt.Println("Usage: ibsc [opt] <command>")
	fmt.Println()
	fmt.Println("IBS client resolves .ibs domains over http/https before running the command.")
	fmt.Println("All strings which end with .ibs will be replaced by their assigned ip.")
	fmt.Println("Server selection is done in the config file located at $HOME/.ibsc_conf")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("%-15s: %s", "-v, --version", "Print command version\n")
	fmt.Printf("%-15s: %s", "-h, --help", "Print usage information\n")
}
