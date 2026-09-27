package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Topic struct {
	Topic string `json:"topic"`
	Source string `json:"source"`
	Name string `json:"name"`
	Description string `json:"description"`
	Unit string `json:"unit"`
}

var topics []Topic

func loadTopics(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("could not read %s: %v", path, err)
	}

	err = json.Unmarshal(data, &topics)
	if err != nil {
		log.Fatalf("could not parse %s: %v", path, err)
	}

	for i := range topics {
		parts := strings.Split(topics[i].Topic, "/")
		if len(parts) == 3 {
			topics[i].Source = parts[1]
		}
	}

	log.Printf("loaded %d topics from %s", len(topics), path)
}

func handleTopics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(topics)
}
