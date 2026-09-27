package main

import (
	"log"
	"net/http"
)

var hub = newHub()

func main() {
	connectMQTT()

	http.HandleFunc("/events", handleEvents)

	log.Println("listening on http://localhost:1880/events")
	log.Fatal(http.ListenAndServe(":1880", nil))
}
