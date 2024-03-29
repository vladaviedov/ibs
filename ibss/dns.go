package main

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

var badTypeError = errors.New("Unsupported record type requested")
var notFoundError = errors.New("DNS record not found")
var notHandledError = errors.New("Not handling this address")

func ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		answer, err := resolve(question.Name, question.Qtype)
		if err != nil {
			log.Println("Failed to resolve record: ", err.Error())
			continue;
		}

		msg.Answer = append(msg.Answer, answer...)
	}

	w.WriteMsg(&msg)
}

func resolve(name string, qtype uint16) ([]dns.RR, error) {
	if qtype != dns.TypeA {
		return nil, badTypeError
	}

	// Format: device.ibs.
	parts := strings.Split(name, ".")
	if len(parts) > 3 || parts[1] != "ibs" {
		return nil, notHandledError
	}

	device, ok := DeviceMap[parts[0]]
	if !ok {
		return nil, notFoundError
	}

	return []dns.RR{&dns.A{
		Hdr: dns.RR_Header{
			Name: name,
			Rrtype: dns.TypeA,
			Class: dns.ClassINET,
			Ttl: 60,
		},
		A: net.ParseIP(device.IP),
	}}, nil
}
