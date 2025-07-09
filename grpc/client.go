package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/biaferre/grpc-sensor-server/sensor/pbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewSensorClient(conn)
	ctx, cancel := context.WithTimeout((context.Background()), 120*time.Second)
	defer cancel()

	var number string

	log.Println("Client started, define a number of sensors: ")
	fmt.Scanln(&number)

	startTime := time.Now()

	r, err := client.GetSensorData(ctx, &pb.SensorRequest{
		RequestMessage: number,
	})

	if err != nil {
		log.Fatalf("Error calling GetSensorData: %v", err)
	}

	var errorMessage string
	if r.Err != "" {
		errorMessage = r.Err
	} else {
		errorMessage = "No error"
	}
	log.Printf("Response from server: AvgTemp: %f, MinTemp: %f, MaxTemp: %f, Err: %s",
		r.AvgTemp, r.MinTemp, r.MaxTemp, errorMessage)

	for _, sensor := range r.Sensors {
		log.Printf("Sensor ID: %d, Temperature: %d", sensor.SensorId, sensor.Temperature)
	}
	log.Println("RTT:", time.Since(startTime).Milliseconds(), "ms")
}
