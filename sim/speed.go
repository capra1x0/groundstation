package main

import (
	"math"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	speedTick = 250 * time.Microsecond
	stopTime = 3 * time.Second
)

type phase struct {
	target float64
	rate float64
}

var driveCycle = []phase {
	{ target: 100, rate: 16 },
	{ target: 60, rate: 10 },
	{ target: 100, rate: 12 },
	{ target: 0, rate: 20 },
}

func runSpeed(client mqtt.Client) {
	speed := 0.0

	for {
		for _, p := range driveCycle {
			step := p.rate * speedTick.Seconds()

			for speed != p.target {
				speed = moveTowards(speed, p.target, step)
				publish(client, "speed", speed, "km/h")
				time.Sleep(speedTick)
			}
		}

		stopEnd := time.Now().Add(stopTime)
		for time.Now().Before(stopEnd) {
			publish(client, "speed", 0, "km/h")
			time.Sleep(speedTick)
		}
	}
}

func moveTowards(current, target, step float64) float64 {
	if current < target {
		return math.Min(current + step, target)
	}
	return math.Max(current - step, target)
}
