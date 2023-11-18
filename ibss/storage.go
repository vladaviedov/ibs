package main

import (
	"net/netip"
	"time"

	"github.com/gorilla/mux"
)

type DeviceDetails struct {
	Address   netip.Addr
	MAC       string
	Timestamp time.Time
}

var Router = mux.NewRouter().StrictSlash(true)
var DeviceMap = make(map[string]DeviceDetails)
var Passkey string
