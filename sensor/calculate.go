package sensor

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	pbs "github.com/biaferre/grpc-sensor-server/sensor/pbs"
	"google.golang.org/protobuf/encoding/protojson"
)

type SensorData struct {
	ID          int32 `json:"sensor_id"`
	Temperature int32 `json:"temperature"`
}

func sensor(wg *sync.WaitGroup, ch chan int, numSensor int, duration time.Duration) int {
	var data int
	time.Sleep(duration)
	defer wg.Done()
	data = rand.Intn(100)
	ch <- data

	fmt.Println("Sensor ", numSensor, ": duration: ", duration, ", temperature:", data)
	return data
}

func startTimer() time.Time {
	fmt.Println("START")
	return time.Now()
}

func endTimer(startTime time.Time) {
	endTime := time.Now()
	fmt.Println("THE END / duration:", endTime.Sub(startTime))
}

func formatResults(sensorData []int) (s []byte) {
	var max = sensorData[0]
	var min = sensorData[0]
	var avg = float32(0)

	for _, data := range sensorData {
		avg += float32(data)
		if data > max {
			max = data
		}
		if data < min {
			min = data
		}
	}

	avg /= float32(len(sensorData))

	var sensors []SensorData
	for i, data := range sensorData {
		sensors = append(sensors, SensorData{
			ID:          int32(i + 1),
			Temperature: int32(data),
		})
	}

	data := map[string]interface{}{
		"avg_temp": avg,
		"max_temp": max,
		"min_temp": min,
		"sensors":  sensors,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println("Results:\n" + string(jsonBytes))

	return jsonBytes
}

func AggregateData(wg *sync.WaitGroup, sensorData []int, durations []time.Duration, numSensors int) (s []byte) {
	defer endTimer(startTimer())

	wg.Add(numSensors)

	sensorDataChan := make(chan int, numSensors)

	for i := 1; i <= numSensors; i++ {
		go sensor(wg, sensorDataChan, i, durations[i-1])
	}

	wg.Wait()
	close(sensorDataChan)

	for value := range sensorDataChan {
		sensorData = append(sensorData, value)
	}

	return formatResults(sensorData)
}

func FormatResponse(rawResponse []byte) *pbs.SensorResponse {
	response := &pbs.SensorResponse{}
	err := protojson.Unmarshal(rawResponse, response)
	if err != nil {
		return &pbs.SensorResponse{Err: err.Error()}
	}

	return response
}
