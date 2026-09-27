package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Topic struct {
	Topic string `json:"topic"`
	Name string `json:"name"`
	Description string `json:"description"`
	SourceLabel string `json:"sourceLabel"`
	Unit string `json:"unit"`
	ValueType string `json:"valueType"`
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

	log.Printf("loaded %d topics from %s", len(topics), path)
}

func handleTopics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(topics)
}
