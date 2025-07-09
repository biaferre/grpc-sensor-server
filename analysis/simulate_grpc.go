package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	pbs "github.com/biaferre/grpc-sensor-server/sensor/pbs"
	"google.golang.org/grpc"
)

func SimulateGrpc() []Result {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pbs.NewSensorClient(conn)

	inputSizes := []int{0, 250, 500, 750, 1000, 1250, 1500, 1750, 2000}
	var results []Result

	for _, size := range inputSizes {
		start := time.Now()

		resp, err := client.GetSensorData(context.Background(), &pbs.SensorRequest{
			RequestMessage: strconv.Itoa(size),
		})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		duration := time.Since(start)
		log.Printf("Input: %d | Latency: %d | Result: %s", size, duration.Milliseconds(), resp)

		results = append(results, Result{InputSize: size, Latency: duration})
	}

	fmt.Println("gRPC benchmark complete. Results saved.")

	return results
}
