package main

import (
	"errors"
	"log"
	"net"

	"github.com/miekg/dns"
)

var BadTypeError = errors.New("Unsupported record type requested")
var NotFoundErorr = errors.New("DNS record not found")

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
		return nil, BadTypeError
	}

	if name == "vladpi.ibs." {
		return []dns.RR{&dns.A{
			Hdr: dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeA,
				Class: dns.ClassINET,
				Ttl: 60,
			},
			A: net.ParseIP("192.168.0.1"),
		}}, nil
	}

	return nil, NotFoundErorr
}
