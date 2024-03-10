package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Report struct {
	Identifier string `json:"identifier"`
	MAC string `json:"mac"`
	IP string `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
	Passkey string `json:"passkey"`
}

func ShowDevices(w http.ResponseWriter, r *http.Request) {
	// Header
	fmt.Fprintf(w, "%-20s %-20s %-20s %-20s %-20s\n\n",
		"Identifier",
		"IP",
		"Client TS",
		"Server TS",
		"Status",
	)

	for _, device := range(DeviceMap) {
		fmt.Fprintf(w, "%-20s %-20s %-20s %-20s %-20s\n",
			device.Identifier,
			device.IP,
			device.ClientTimestamp.Format(time.DateTime),
			device.ServerTimestamp.Format(time.DateTime),
			getStatus(device.ServerTimestamp),
		)
	}
}

func ProcessReport(w http.ResponseWriter, r *http.Request) {
	// Parse report
	var report Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Bad argument")
		return
	}

	// Check auth
	if !checkAuth(&report) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Permission denied")
		return
	}

	// Check storage
	DeviceMutex.Lock()
	store, exists := DeviceMap[report.Identifier]
	if exists {
		if store.MAC != report.MAC {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Occupied")
			return
		}
	}

	DeviceMap[report.Identifier] = DeviceData{
		Identifier: report.Identifier,
		MAC: report.MAC,
		IP: report.IP,
		ClientTimestamp: report.Timestamp.UTC(),
		ServerTimestamp: time.Now().UTC(),
	}
	DeviceMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func ResolveOverHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	identifier, ok := vars["id"]

	device, ok := DeviceMap[identifier]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Device not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, device.IP)
}

// TODO: better auth system
func checkAuth(report *Report) bool {
	return report.Passkey == Config.Passkey
}

func getStatus(lastTimestamp time.Time) string {
	minElapsed := time.Now().Sub(lastTimestamp).Minutes()

	switch {
		case minElapsed <= 5.1:
			return "online"
		case minElapsed <= 30.1:
			return "missing"
		default:
			return "offline"
	}
}
