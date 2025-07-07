package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/biaferre/grpc-sensor-server/sensor"
	amqp "github.com/rabbitmq/amqp091-go"
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

	q, err := ch.QueueDeclare(
		"sensor_data",
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

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag
		true,   // auto-acknowledge
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			log.Panicln(err, "Failed to convert body to integer")

			log.Printf(" [.] fib(%d)", n)
			response := sensor.GenerateSensorData(&n)

			err = ch.PublishWithContext(ctx,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			log.Panicln(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	go func() {
		for d := range msgs {
			bodyStr := string(d.Body)
			val, err := strconv.Atoi(bodyStr)
			if err != nil {
				log.Printf("Error converting message body to int: %v", err)
				continue
			}

			sensor.GenerateSensorData(&val)
		}
	}()

	log.Println("Successfully connected to RabbitMQ and waiting for messages...")
	<-forever
}
