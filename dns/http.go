package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Report struct {
	DeviceData
	Passkey string `json:"passkey"`
}

func ShowDevices(w http.ResponseWriter, r *http.Request) {
	for _, device := range(DeviceMap) {
		log.Println(device.Timestamp)
		fmt.Fprintf(w, "%-20s %-15s %s\n",
			device.Identifier,
			device.IP,
			device.Timestamp.Format(time.DateTime),
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

	DeviceMap[report.Identifier] = report.DeviceData
	DeviceMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func checkAuth(report *Report) bool {
	return report.Passkey == "test"
}