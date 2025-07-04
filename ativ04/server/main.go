package main

import (
	"context"
	"log"
	"net"

	"github.com/biaferre/grpc-sensor-server/sensor"
	"google.golang.org/grpc"
)

type SensorServer struct {
	sensor.UnimplementedSensorServer
}

func (s SensorServer) GetSensorData(ctx context.Context, req *sensor.SensorRequest) (*sensor.SensorResponse, error) {
	log.Printf("Received request for sensor data")
	return &sensor.SensorResponse{
		AvgTemp: 40.0,
		MinTemp: 20.0,
		MaxTemp: 60.0,
		Err:     "",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	serviceRegistrar := grpc.NewServer()
	service := &SensorServer{}
	sensor.RegisterSensorServer(serviceRegistrar, service)

	err = serviceRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}

	log.Println("Server started on localhost:8080")
}
