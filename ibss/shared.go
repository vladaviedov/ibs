package main

import (
	"sync"
	"time"
)

type Settings struct {
	HTTP struct {
		Use bool `json:"use"`
		Port string `json:"port"`
		DNSResolver bool `json:"dnsResolver"`
	} `json:"http"`
	HTTPS struct {
		Use bool `json:"use"`
		Port string `json:"port"`
		DNSResolver bool `json:"dnsResolver"`
		CertFile string `json:"certFile"`
		KeyFile string `json:"keyFile"`
	} `json:"https"`
	DNS struct {
		Use bool `json:"use"`
		Port string `json:"port"`
	} `json:"dns"`
	Passkey string `json:"passkey"`
}

type DeviceData struct {
	Identifier string
	MAC string
	IP string
	ClientTimestamp time.Time
	ServerTimestamp time.Time
}

var Config Settings
var DeviceMap = make(map[string]DeviceData)
var DeviceMutex sync.Mutex
