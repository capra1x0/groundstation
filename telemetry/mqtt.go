package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Reading struct {
	Source string `json:"source"`
	Sensor string `json:"sensor"`
	TS int64 `json:"ts"`
	Value float64 `json:"value"`
	Unit string `json:"unit"`
}

func connectMQTT() {
	options := mqtt.NewClientOptions()
	options.AddBroker("tcp://localhost:1883")
	options.SetConnectRetry(true)
	options.SetConnectRetryInterval(2 * time.Second)
	options.SetAutoReconnect(true)

	options.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("connected to broker")
		c.Subscribe("telemetry/#", 0, onMessage)
	})

	client := mqtt.NewClient(options)

	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func onMessage(_ mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) != 3 {
		return
	}

	var reading Reading
	err := json.Unmarshal(msg.Payload(), &reading)
	if err != nil {
		log.Printf("bad message on %s: %v", msg.Topic(), err)
		return
	}

	reading.Source = parts[1]
	reading.Sensor = parts[2]

	hub.broadcase(reading)
}
