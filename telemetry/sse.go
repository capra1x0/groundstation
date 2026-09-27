package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	readings := hub.add()
	defer hub.remove(readings)

	for {
		select {
		case <- r.Context().Done():
			return

		case reading := <- readings:
			data, _ := json.Marshal(reading)
			fmt.Fprintf(w, "data: %s\n\n", data)
			http.NewResponseController(w).Flush()
		}
	}
}
