package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Status struct {
	DeviceName string    `json:"deviceName"`
	Address    string    `json:"address"`
	Timestamp  time.Time `json:"timestamp"`
}

func DeviceStatus(w http.ResponseWriter, r *http.Request) {
	deviceName := r.URL.String()[1:]
	deviceRecord, ok := DeviceMap[deviceName]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal error")
		return
	}

	status := Status{
		DeviceName: deviceName,
		Address:    deviceRecord.Address.String(),
		Timestamp:  deviceRecord.Timestamp,
	}
	json.NewEncoder(w).Encode(status)
}

func ListDevices(w http.ResponseWriter, r *http.Request) {
	var list []Status
	currentTime := time.Now()

	for device, record := range DeviceMap {
		old := record.Timestamp.Add(10 * time.Minute).Before(currentTime)
		if old {
			continue
		}

		list = append(list, Status{
			DeviceName: device,
			Address:    record.Address.String(),
			Timestamp:  record.Timestamp,
		})
	}

	json.NewEncoder(w).Encode(list)
}
