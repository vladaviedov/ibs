package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"time"
)

type Report struct {
	DeviceName    string    `json:"deviceName"`
	DeviceMAC     string    `json:"deviceMac"`
	AddressReport string    `json:"addressReport"`
	Timestamp     time.Time `json:"timestamp"`
	Passkey       string    `json:"passkey"`
}

func ProcessReport(w http.ResponseWriter, r *http.Request) {
	// Parse report body
	var report Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Bad arguments")
		return
	}

	// Access device map
	var detailsPtr *DeviceDetails = nil
	details, ok := DeviceMap[report.DeviceName]
	if !ok {
		Router.HandleFunc(fmt.Sprintf("/%s", report.DeviceName), DeviceStatus)
	} else {
		detailsPtr = &details
	}

	// Check allowed
	if !isAllowed(detailsPtr, &report) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Permission denied")
		return
	}

	// Parse ip
	ip, err := netip.ParseAddr(report.AddressReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid IP address")
		return
	}

	DeviceMap[report.DeviceName] = DeviceDetails{
		Address:   ip,
		MAC:       report.DeviceMAC,
		Timestamp: report.Timestamp,
	}
}

func isAllowed(current *DeviceDetails, report *Report) bool {
	if report.Passkey != Passkey {
		return false
	}

	if current != nil {
		return current.MAC == report.DeviceMAC
	}

	return true
}
