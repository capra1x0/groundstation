package main

import (
	"math/rand/v2"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func runTemp(client mqtt.Client) {
	temp := 23.0	

	for {
		random := rand.IntN(10)
		randChange := rand.IntN(20)
		if random <= 1 {
			temp -= float64(randChange) / 10
		}
		if random >= 8 {
			temp += float64(randChange) / 10
		}

		publish(client, "temp", temp, "Celcius")
		time.Sleep(time.Second)
	}
}
