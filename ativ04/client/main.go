package main

import (
	"context"
	"log"
	"time"

	pb "github.com/biaferre/grpc-sensor-server/sensor"
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
	ctx, cancel := context.WithTimeout((context.Background()), 2*time.Second)
	defer cancel()

	startTime := time.Now()
	r, err := client.GetSensorData(ctx, &pb.SensorRequest{
		RequestMessage: "Hey its bia",
	})
	endTime := time.Now()

	if err != nil {
		log.Fatalf("Error calling GetSensorData: %v", err)
	}
	log.Printf("Response from server: AvgTemp: %f, MinTemp: %f, MaxTemp: %f, Err: %s",
		r.AvgTemp, r.MinTemp, r.MaxTemp, r.Err)
	log.Println("Client finished successfully in:", endTime.Sub(startTime).Milliseconds(), "ms")
}
