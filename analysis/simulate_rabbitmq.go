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

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}
	return string(result)
}

func SimulateRabbit() []Result {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare reply queue: %v", err)
	}

	msgs, err := ch.Consume(
		replyQueue.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume from reply queue: %v", err)
	}

	inputSizes := []int{0, 250, 500, 750, 1000, 1250, 1500, 1750, 2000}
	var results []Result

	for _, size := range inputSizes {
		corrID := randomString(16)
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := ch.PublishWithContext(ctx,
			"",
			"sensor_data",
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrID,
				ReplyTo:       replyQueue.Name,
				Body:          []byte(fmt.Sprintf("%d", size)),
			})
		if err != nil {
			log.Fatalf("Failed to publish: %v", err)
		}

		for d := range msgs {
			if d.CorrelationId == corrID {
				latency := time.Since(start)
				log.Printf("Input: %d | Latency: %s | Response: %s", size, latency, string(d.Body))
				results = append(results, Result{InputSize: size, Latency: latency})
				break
			}
		}
	}

	fmt.Println("RabbitMQ benchmark complete. Results saved.")

	return results
}
