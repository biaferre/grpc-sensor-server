package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Result struct {
	InputSize int
	Latency   time.Duration
}

func writeResultsToCSV(results []Result, title string) {
	file, err := os.Create(title + ".csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"InputSize", "LatencyMicroseconds"})
	for _, r := range results {
		writer.Write([]string{
			strconv.Itoa(r.InputSize),
			strconv.FormatInt(r.Latency.Microseconds(), 10),
		})
	}
	fmt.Println("Results saved to grpc_results.csv")
}

func main() {
	grpcResults := SimulateGrpc()
	writeResultsToCSV(grpcResults, "grpc_results")

	fmt.Println("Data saved to grpc_results.csv")

	rabbitResults := SimulateRabbit()
	writeResultsToCSV(rabbitResults, "rabbitmq_results")

	fmt.Println("Data saved to rabbitmq_results.csv")
}
