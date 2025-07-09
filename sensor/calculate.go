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

func sensor(wg *sync.WaitGroup, ch chan SensorData, numSensor int, duration time.Duration) SensorData {
	time.Sleep(duration)
	defer wg.Done()

	data := SensorData{
		ID:          int32(numSensor),
		Temperature: int32(rand.Intn(100)),
	}
	ch <- data

	fmt.Println("Sensor ", numSensor, ": duration: ", duration, ", temperature:", data.Temperature)
	return data
}

func AggregateData(wg *sync.WaitGroup, sensorData []SensorData, durations []time.Duration, numSensors int) (s []byte) {
	wg.Add(numSensors)

	sensorDataChan := make(chan SensorData, numSensors)

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

func formatResults(sensorData []SensorData) (s []byte) {
	var max = sensorData[0].Temperature
	var min = sensorData[0].Temperature
	var avg = float32(0)

	for _, data := range sensorData {
		avg += float32(data.Temperature)
		if data.Temperature > max {
			max = data.Temperature
		}
		if data.Temperature < min {
			min = data.Temperature
		}
	}

	avg /= float32(len(sensorData))

	var sensors []SensorData
	for _, data := range sensorData {
		sensors = append(sensors, SensorData{
			ID:          int32(data.ID),
			Temperature: int32(data.Temperature),
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

func FormatResponse(rawResponse []byte) *pbs.SensorResponse {
	response := &pbs.SensorResponse{}
	err := protojson.Unmarshal(rawResponse, response)
	if err != nil {
		return &pbs.SensorResponse{Err: err.Error()}
	}

	return response
}
