package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

func main() {
	dnsServ := &dns.Server{
		Addr: ":8080",
		Net: "udp",
		UDPSize: 25565,
		Handler: dns.HandlerFunc(ServeDNS),
		ReusePort: true,
	}

	httpRouter := mux.NewRouter().StrictSlash(true)
	httpRouter.HandleFunc("/", ShowDevices)

	err := dnsServ.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server: ", err.Error())
	}
}

