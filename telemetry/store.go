package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

var pending = make(chan Reading, 10000)

func connectDatabase() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("invalid DATABASE_URL: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("could not reach databse: %v", err)
	}
	
	db = pool
	log.Println("connected to database")

	go writeLoop()
}

func store(reading Reading) {
	select {
	case pending <- reading:
	default:
		log.Println("database writer is behind, dropping reading")
	}
}

func writeLoop() {
	ticker := time.NewTicker(time.Second)
	batch := make([]Reading, 0, 500)

	for {
		select {
		case reading := <-pending:
			batch = append(batch, reading)
			if len(batch) >= 500 {
				batch = saveBatch(batch)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				batch = saveBatch(batch)
			}
		}
	}
}

func saveBatch(batch []Reading) []Reading {
	rows := make([][]any, 0, len(batch))
	for _, reading := range batch {
		row := toRow(reading)
		if row != nil {
			rows = append(rows, row)
		}	
	}

	_, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"readings"},
		[]string{"ts", "topic", "value_num", "value_bool", "value_text"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Printf("saving %d readings failed: %v", len(rows), err)
	}

	return batch[:0]
}

func toRow(reading Reading) []any {
	ts := time.UnixMilli(reading.TS)
	if reading.TS == 0 {
		ts = time.Now()
	}

	switch value := reading.Value.(type) {
	case float64:
		return []any{ts, reading.Topic, value, nil, nil}
	case bool:
		return []any{ts, reading.Topic, nil, value, nil}
	case string:
		return []any{ts, reading.Topic, nil, nil, value}
	default:
		return nil
	}
}
