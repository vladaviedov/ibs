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
const version = "0.3.0"

type Settings struct {
	Server string `json:"server"`
	PortHTTP uint `json:"portHttp"`
	PortHTTPS uint `json:"portHttps"`
	WithShell bool `json:"withShell"`
	ShowCommand bool `json:"showCommand"`
}

func main() {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad matching pattern")
		fmt.Fprintln(os.Stderr, "Your ibsc binary is broken :(")
		os.Exit(1)
	}

	cmdStartIndex := len(os.Args)
	var configPath *string = nil

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if (arg[0] != '-') {
			cmdStartIndex = i
			break
		}

		if arg == "--version" || arg== "-v" {
			printVersion()
			os.Exit(0)
		}
		if arg == "--help" || arg == "-h" {
			printUsage()
			os.Exit(0)
		}
		if arg == "--config" || arg == "-c" {
			// Check if parameter exists
			if len(os.Args) <= i + 1 || configPath != nil {
				fmt.Fprintln(os.Stderr, "Invalid config option invokation")
				os.Exit(1)
			}

			configPath = &os.Args[i + 1]
			i++
		}
	}

	config := loadConfig(configPath)
	if cmdStartIndex == len(os.Args) {
		printStatus(&config)
		os.Exit(0)
	}
	
	// Find ibs addresses
	var argv []string
	for _, arg := range os.Args[cmdStartIndex:] {
		for {
			match := reg.FindString(arg)
			if match == "" {
				break
			}

			ip := resolveDomain(&config, match)
			arg = strings.ReplaceAll(arg, match, ip)
		}

		argv = append(argv, arg)
	}

	var bin string
	if !config.WithShell {
		// Find needed binary, since we don't have execvp
		bin, err = exec.LookPath(argv[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
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
		argv = append([]string{bin, "-c"}, strings.Join(argv, " "))
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

func loadConfig(path *string) Settings {
	var configPath string
	if path == nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "No home directory :(")
			os.Exit(1)
		}

		// Default config path
		configPath = homeDir + "/.ibsc_conf"
	} else {
		configPath = *path
	}

	// Read config
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

func printStatus(config *Settings) {
	fmt.Println("Server:", config.Server)
	fmt.Println("--------------------")

	// Ping HTTPS
	url := fmt.Sprintf("https://%s:%d/ping", config.Server, config.PortHTTPS)
	res, err := http.Get(url)
	fmt.Printf("https (%d): ", config.PortHTTPS)
	if err != nil {
		fmt.Println("not reachable")
	} else if res.StatusCode != http.StatusOK {
		fmt.Println("server error")
	} else {
		fmt.Println("good")
	}

	// Ping HTTP
	url = fmt.Sprintf("http://%s:%d/ping", config.Server, config.PortHTTP)
	res, err = http.Get(url)
	fmt.Printf("http (%d): ", config.PortHTTP)
	if err != nil {
		fmt.Println("not reachable")
	} else if res.StatusCode != http.StatusOK {
		fmt.Println("server error")
	} else {
		fmt.Println("good")
	}
}

func resolveDomain(config *Settings, domain string) string {
	identifier := strings.TrimSuffix(domain, ".ibs")

	// Try https
	url := fmt.Sprintf("https://%s:%d/dns/%s", config.Server, config.PortHTTPS, identifier)
	res, err := http.Get(url)
	if err != nil {
		// Try http
		url = fmt.Sprintf("http://%s:%d/dns/%s", config.Server, config.PortHTTP, identifier)
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
	fmt.Printf("%-20s: %s\n", "-v, --version", "Print command version")
	fmt.Printf("%-20s: %s\n", "-h, --help", "Print usage information")
	fmt.Printf("%-20s: %s\n", "-c, --config <path>", "Specify configuration file")
}
