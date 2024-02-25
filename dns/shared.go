package main

import (
	"sync"
	"time"
)

type DeviceData struct {
	Identifier string `json:"identifier"`
	MAC string `json:"mac"`
	IP string `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

var DeviceMap = make(map[string]DeviceData)
var DeviceMutex sync.Mutex
