package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func randomString(l int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, l)
	for i := range result {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[nBig.Int64()]
	}
	return string(result)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to declare reply queue: %v", err)
	}

	msgs, err := ch.Consume(
		replyQueue.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Client started, define a number of sensors: ")
	var number int
	fmt.Scanln(&number)

	corrID := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	startTime := time.Now()

	err = ch.PublishWithContext(ctx,
		"",            // exchange
		"sensor_data", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       replyQueue.Name,
			Body:          []byte(fmt.Sprintf("%d", number)),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	for d := range msgs {
		if d.CorrelationId == corrID {
			log.Printf("Response: "+string(d.Body)+"\n RTT: %d ms",
				time.Since(startTime).Milliseconds())
			break
		}
	}
}
