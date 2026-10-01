package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

var hub = newHub()

func main() {
	err := godotenv.Load("../.env")		
	if err != nil {
		log.Println("no ../.env found")
	}

	loadTopics("topics.json")
	connectDatabase()
	connectMQTT()

	http.HandleFunc("/events", handleEvents)
	http.HandleFunc("/topics", handleTopics)

	log.Println("listening on http://localhost:1880/events")
	log.Fatal(http.ListenAndServe(":1880", nil))
}
