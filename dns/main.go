package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

func main() {
	go httpServer()
	dnsServer()
}

func httpServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ShowDevices).Methods(http.MethodGet)
	router.HandleFunc("/", ProcessReport).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func dnsServer() {
	server := &dns.Server{
		Addr: ":8081",
		Net: "udp",
		UDPSize: 25565,
		Handler: dns.HandlerFunc(ServeDNS),
		ReusePort: true,
	}

	log.Fatal(server.ListenAndServe())
}
