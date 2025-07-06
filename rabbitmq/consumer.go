package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		panic(err)
	}

	msgs, err := ch.Consume(
		"sensor_data", // queue name
		"",            // consumer tag
		true,          // auto-acknowledge
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// todo processing
			fmt.Printf("Processing sensor data: %s\n", d.Body)
		}
	}()

	log.Println("Successfully connected to RabbitMQ and waiting for messages...")
	<-forever
}
