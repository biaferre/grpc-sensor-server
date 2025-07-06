package sensor

import (
	"math/rand"
	"sync"
	"time"
)

func GenerateSensorData(numSensors *int) (s []byte) {
	var wg sync.WaitGroup
	var sensorData []int

	var durations []time.Duration
	var ns int

	if numSensors == nil {
		ns = rand.Intn(6)
	} else {
		ns = *numSensors
	}

	for i := 0; i < ns; i++ {
		durations = append(durations, time.Duration(rand.Intn(10))*time.Second)
	}

	return AggregateData(&wg, sensorData, durations, ns)
}
