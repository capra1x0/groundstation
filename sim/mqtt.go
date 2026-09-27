package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type reading struct {
	Ts    int64   `json:"ts"`
	Value any `json:"value"`
	Unit  string  `json:"unit"`
}

func connect(broker string) mqtt.Client {
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)

	options.SetClientID(fmt.Sprintf("groundstation-sim-%d", time.Now().UnixNano()))

	options.SetConnectRetry(true)
	options.SetConnectRetryInterval(2 * time.Second)
	options.SetAutoReconnect(true)

	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("connection lost: %v", err)
	})

	client := mqtt.NewClient(options)

	log.Printf("connecting to %s ...", broker)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		log.Fatalf("connect failed: %v", token.Error())
	}

	return client
}

func publish(client mqtt.Client, sensor string, value any, unit string) {
	if !client.IsConnectionOpen() {
		return
	}

	payload, err := json.Marshal(reading{
		Ts: time.Now().UnixMilli(),
		Value: value,
		Unit: unit,
	})
	if err != nil {
		log.Printf("could not encode %s: %v", sensor, err)
		return
	}

	client.Publish("telemetry/sim/"+sensor, 0, false, payload)
}
