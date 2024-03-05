package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

const version = "0.1.0"

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
		configPath = "config.json"
	}

	loadSettings(configPath)

	if Config.HTTP.Use {
		go httpServer()
	}
	if Config.HTTPS.Use {
		go httpsServer()
	}

	dnsServer()
}

func loadSettings(configPath string) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to open " + configPath + ". A config file is required")
	}

	err = json.Unmarshal(content, &Config)
	if err != nil {
		log.Fatal("Failed to parse config file")
	}

	log.Println("Server config loaded")
}

func httpServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ShowDevices).Methods(http.MethodGet)
	router.HandleFunc("/", ProcessReport).Methods(http.MethodPost)
	router.HandleFunc("/ping", ping).Methods(http.MethodGet)

	if Config.HTTP.DNSResolver {
		router.HandleFunc("/dns/{id}", ResolveOverHTTP).Methods(http.MethodGet)
	}

	log.Fatal(http.ListenAndServe(Config.HTTP.Port, router))
}

func httpsServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ShowDevices).Methods(http.MethodGet)
	router.HandleFunc("/", ProcessReport).Methods(http.MethodPost)
	router.HandleFunc("/ping", ping).Methods(http.MethodGet)

	if Config.HTTPS.DNSResolver {
		router.HandleFunc("/dns/{id}", ResolveOverHTTP).Methods(http.MethodGet)
	}

	log.Fatal(http.ListenAndServeTLS(
		Config.HTTPS.Port,
		Config.HTTPS.CertFile,
		Config.HTTPS.KeyFile,
		router,
	))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func dnsServer() {
	server := &dns.Server{
		Addr: Config.DNS.Port,
		Net: "udp",
		UDPSize: 25565,
		Handler: dns.HandlerFunc(ServeDNS),
		ReusePort: true,
	}

	log.Fatal(server.ListenAndServe())
}

func printVersion() {
	fmt.Printf("IBS server v%s\n", version)
}

func printUsage() {
	fmt.Println("Usage: ibss [opt] [configpath]")
	fmt.Println()
	fmt.Println("IBS server is a DNS/HTTP/HTTPS server for the IBS system.")
	fmt.Println("Please see the config file to configure your server components.")
	fmt.Println("By default ./config.json will be parsed for settings.")
	fmt.Println("The used config file path can be changed by passing a filepath as an argument.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("%-15s: %s", "-v, --version", "Print command version\n")
	fmt.Printf("%-15s: %s", "-h, --help", "Print usage information\n")
}
