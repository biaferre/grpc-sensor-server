package main

import (
	"fmt"
	"log"

	"github.com/biaferre/grpc-sensor-server/sensor"
	"github.com/streadway/amqp"
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
		"sensor_data",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		panic(err)
	}

	fmt.Println(q)

	body := sensor.GenerateSensorData(nil)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		panic(err)
	}

	log.Printf("Message sent to queue: %s", q.Name)
}
