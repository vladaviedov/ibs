package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeviceData struct {
	Identifier string `json:"identifier"`
	MAC string `json:"mac"`
	IP string `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

type Report struct {
	DeviceData
	Passkey string `json:"passkey"`
}

var deviceMap = make(map[string]DeviceData)

func ShowDevices(w http.ResponseWriter, r *http.Request) {
	for _, device := range(deviceMap) {
		fmt.Fprintf(w, "%20s %15s %s\n",
			device.Identifier,
			device.IP,
			device.Timestamp.Format(time.DateTime),
		)
	}

	w.WriteHeader(http.StatusOK)
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
	store, exists := deviceMap[report.Identifier]
	if exists {
		if store.MAC != report.MAC {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Occupied")
			return
		}
	}

	deviceMap[report.Identifier] = report.DeviceData
	w.WriteHeader(http.StatusOK)
}

func checkAuth(report *Report) bool {
	return report.Passkey == "test"
}
