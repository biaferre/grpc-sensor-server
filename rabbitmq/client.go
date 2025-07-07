package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		panic(err)
	}
	defer conn.Close()

	log.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		panic(err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
		panic(err)
	}

	fmt.Println(q)
	log.Println("Client started, define a number of sensors: ")

	var number int
	fmt.Scanln(&number)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",            // exchange
		"sensor_data", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%d", number)),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		panic(err)
	}

	// Start consuming messages
	go func() {
		for d := range msgs {
			bodyStr := string(d.Body)
			log.Printf("Received message: %s", bodyStr)
		}
	}()

	log.Printf("Message sent to queue: %s", q.Name)
}
