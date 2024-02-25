package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

func main() {
	loadSettings()

	if Config.HTTP.Use {
		go httpServer()
	}
	if Config.HTTPS.Use {
		go httpsServer()
	}

	dnsServer()
}

func loadSettings() {
	content, err := os.ReadFile("settings.json")
	if err != nil {
		log.Fatal("Failed to open 'settings.json'. Config file is required")
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

	if Config.HTTP.DNSResolver {
		router.HandleFunc("/dns/{id}", ResolveOverHTTP).Methods(http.MethodGet)
	}

	log.Fatal(http.ListenAndServe(Config.HTTP.Port, router))
}

func httpsServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ShowDevices).Methods(http.MethodGet)
	router.HandleFunc("/", ProcessReport).Methods(http.MethodPost)

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
