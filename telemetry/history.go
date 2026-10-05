package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const defaultHistoryRange = 5 * time.Minute

const historyQuery = `
	SELECT ts, value_num, value_bool, value_text
	FROM readings
	WHERE topic = $1 AND ts >= $2 AND ts <= $3
	ORDER BY ts
`

func handleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query()

	topic := query.Get("topic")
	if topic == "" {
		http.Error(w, "topic is required", http.StatusBadRequest)
		return
	}

	to, ok := parseMillis(query.Get("to"), time.Now())
	if !ok {
		http.Error(w, "invalid value for to", http.StatusBadRequest)
		return
	} 

	from, ok := parseMillis(query.Get("from"), to.Add(-defaultHistoryRange))
	if !ok {
		http.Error(w, "invalid value for from", http.StatusBadRequest)
		return
	}

	if !from.Before(to) {
		http.Error(w, "from must be before to", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(r.Context(), historyQuery, topic, from, to)
	if err != nil {
		log.Printf("history query faied: %v", err)
		http.Error(w, "could not load history ", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	readings := []Reading{}

	for rows.Next() {
		var ts time.Time
		var number *float64
		var boolean *bool
		var text *string

		err := rows.Scan(&ts, &number, &boolean, &text)
		if err != nil {
			log.Printf("reading history row failed: %v", err)
			http.Error(w, "could not load history", http.StatusInternalServerError)
			return
		}

		reading := Reading{Topic: topic, TS: ts.UnixMilli()}

		switch {
		case number != nil:
			reading.Value = *number
		case boolean != nil:
			reading.Value = *boolean
		case text != nil:
			reading.Value = *text
		}

		readings = append(readings, reading)
	}

	if rows.Err() != nil {
		log.Printf("history query failed: %v", rows.Err())
		http.Error(w, "could not load history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(readings)
}

func parseMillis(value string, fallback time.Time) (time.Time, bool) {
	if value == "" {
		return fallback, true
	}

	ms, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, false
	}

	return time.UnixMilli(ms), true
}