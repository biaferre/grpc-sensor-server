package main

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/biaferre/grpc-sensor-server/sensor"

	pbs "github.com/biaferre/grpc-sensor-server/sensor/pbs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type SensorServer struct {
	pbs.UnimplementedSensorServer
}

func (s SensorServer) GetSensorData(ctx context.Context, req *pbs.SensorRequest) (*pbs.SensorResponse, error) {
	log.Printf("Received request for sensor data with number of sensors = %s", req.RequestMessage)

	numSensors, err := strconv.Atoi(req.RequestMessage)
	if err != nil {
		return nil, err
	}

	var rawResponse = sensor.GenerateSensorData(&numSensors)
	return formatResponse(rawResponse), nil
}

func formatResponse(rawResponse []byte) *pbs.SensorResponse {
	response := &pbs.SensorResponse{}
	err := protojson.Unmarshal(rawResponse, response)
	if err != nil {
		return &pbs.SensorResponse{Err: err.Error()}
	}

	return response
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	serviceRegistrar := grpc.NewServer()
	service := &SensorServer{}
	pbs.RegisterSensorServer(serviceRegistrar, service)

	err = serviceRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}

	log.Println("Server started on localhost:8080")
}
