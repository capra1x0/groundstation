package main

import (
	"math/rand/v2"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func runDoor(client mqtt.Client) {
	open := false

	for {
		if rand.IntN(10) == 0 {
			open = !open
		}

		publish(client, "door", open, "")
		time.Sleep(time.Second)
	}
}