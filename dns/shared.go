package main

import (
	"sync"
	"time"
)

type Settings struct {
	HTTP struct {
		Use bool `json:"use"`
		Port string `json:"port"`
	} `json:"http"`
	HTTPS struct {
		Use bool `json:"use"`
		Port string `json:"port"`
		CertFile string `json:"certFile"`
		KeyFile string `json:"keyFile"`
	} `json:"https"`
	DNS struct {
		Port string `json:"port"`
	} `json:"dns"`
}

type DeviceData struct {
	Identifier string `json:"identifier"`
	MAC string `json:"mac"`
	IP string `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

var Config Settings
var DeviceMap = make(map[string]DeviceData)
var DeviceMutex sync.Mutex
