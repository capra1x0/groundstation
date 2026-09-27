package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker adress")
	temp := flag.Bool("temp", false, "simulate temperature")
	speed := flag.Bool("speed", false, "simulate speed")
	door := flag.Bool("door", false, "simulate a door (open/closed)")
	flag.Parse()

	if !*temp && !*speed && !*door {
		fmt.Println("no sensor selected, example: go run . -temp -speed")
		flag.PrintDefaults()
		os.Exit(1)
	}

	client := connect(*broker)
	defer client.Disconnect(250)

	if *temp {
		go runTemp(client)
		log.Println("start temp simulator")
	}
	if *speed {
		go runSpeed(client)
		log.Println("start speed simulator")
	}
	if *door {
		go runDoor(client)
		log.Println("start door simulator")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
}
