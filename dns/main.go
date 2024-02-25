package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

func main() {
	// Make DNS server
	dnsServ := &dns.Server{
		Addr: ":8080",
		Net: "udp",
		UDPSize: 25565,
		Handler: dns.HandlerFunc(ServeDNS),
		ReusePort: true,
	}

	// Make HTTP server
	httpRouter := mux.NewRouter().StrictSlash(true)
	httpRouter.HandleFunc("/", ShowDevices).Methods(http.MethodGet)
	httpRouter.HandleFunc("/", ProcessReport).Methods(http.MethodPost)

	err := dnsServ.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server: ", err.Error())
	}
}

