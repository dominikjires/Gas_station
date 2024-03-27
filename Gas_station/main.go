package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

// Config struct represents the configuration of the simulation.
type Config struct {
	Cars      CarConfig      `yaml:"cars"`
	Stations  StationsConfig `yaml:"stations"`
	Registers RegisterConfig `yaml:"registers"`
}

// CarConfig represents the configuration for cars.
type CarConfig struct {
	Count          int           `yaml:"count"`
	ArrivalTimeMin time.Duration `yaml:"arrival_time_min"`
	ArrivalTimeMax time.Duration `yaml:"arrival_time_max"`
}

// StationsConfig represents the configuration for stations.
type StationsConfig map[string]struct {
	Count        int           `yaml:"count"`
	ServeTimeMin time.Duration `yaml:"serve_time_min"`
	ServeTimeMax time.Duration `yaml:"serve_time_max"`
}

// RegisterConfig represents the configuration for registers.
type RegisterConfig struct {
	Count         int           `yaml:"count"`
	HandleTimeMin time.Duration `yaml:"handle_time_min"`
	HandleTimeMax time.Duration `yaml:"handle_time_max"`
}

// Car represents a car with an ID and arrival time.
type Car struct {
	id          int
	arrivalTime time.Time
}

// Station represents a refueling station.
type Station struct {
	name         string
	count        int
	totalCars    int
	queue        chan struct{}
	serveTimeMin time.Duration
	serveTimeMax time.Duration
	totalTime    time.Duration
	avgQueueTime time.Duration
	maxQueueTime time.Duration
}

// Register represents a payment register.
type Register struct {
	count         int
	totalCars     int
	queue         chan struct{}
	handleTimeMin time.Duration
	handleTimeMax time.Duration
	totalTime     time.Duration
	avgQueueTime  time.Duration
	maxQueueTime  time.Duration
}

// carRefueling simulates car refueling process at a station.
func (station *Station) carRefueling(wg *sync.WaitGroup) {
	defer wg.Done()

	for range station.queue {
		start := time.Now()
		serveTime := station.serveTimeMin + time.Duration(rand.Int63n(int64(station.serveTimeMax-station.serveTimeMin)))
		time.Sleep(serveTime)
		station.totalTime += serveTime
		elapsed := time.Since(start)

		if elapsed > station.maxQueueTime {
			station.maxQueueTime = elapsed
		}
		station.totalCars++
	}
}

// payment simulates payment process at a register.
func (register *Register) payment(wg *sync.WaitGroup) {
	defer wg.Done()

	for range register.queue {
		start := time.Now()
		handleTime := register.handleTimeMin + time.Duration(rand.Int63n(int64(register.handleTimeMax-register.handleTimeMin)))
		time.Sleep(handleTime)
		register.totalTime += handleTime
		elapsed := time.Since(start)

		if elapsed > register.maxQueueTime {
			register.maxQueueTime = elapsed
		}
		register.totalCars++
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Read configuration from file.
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	// Initialize stations and registers.
	stations := make(map[string]*Station)
	registers := &Register{
		count:         config.Registers.Count,
		handleTimeMin: config.Registers.HandleTimeMin,
		handleTimeMax: config.Registers.HandleTimeMax,
		queue:         make(chan struct{}, config.Cars.Count),
	}

	for name, stationConfig := range config.Stations {
		stations[name] = &Station{
			name:         name,
			count:        stationConfig.Count,
			serveTimeMin: stationConfig.ServeTimeMin,
			serveTimeMax: stationConfig.ServeTimeMax,
			queue:        make(chan struct{}, config.Cars.Count),
		}
	}

	// Start goroutines for simulating car refueling at stations.
	var wg sync.WaitGroup
	for _, station := range stations {
		wg.Add(1)
		go station.carRefueling(&wg)
	}

	// Start goroutines for simulating payment at registers.
	for i := 0; i < registers.count; i++ {
		wg.Add(1)
		go registers.payment(&wg)
	}

	// Simulate cars arriving at stations and registers.
	go func() {
		for i := 0; i < config.Cars.Count; i++ {
			time.Sleep(config.Cars.ArrivalTimeMin + time.Duration(rand.Int63n(int64(config.Cars.ArrivalTimeMax-config.Cars.ArrivalTimeMin))))
			station := stations[chooseRandomStation()]
			station.queue <- struct{}{}
			registers.queue <- struct{}{}
		}
		close(registers.queue)
		for _, station := range stations {
			close(station.queue)
		}
	}()

	// Wait for all goroutines to finish.
	wg.Wait()

	// Print statistics for stations.
	fmt.Println("Stations:")
	for name, station := range stations {
		avgQueueTime := time.Duration(0)
		if station.totalCars > 0 {
			avgQueueTime = station.totalTime / time.Duration(station.totalCars)
		}
		fmt.Printf("%s:\n", name)
		fmt.Printf("  total_cars: %d\n", station.totalCars)
		fmt.Printf("  total_time: %s\n", station.totalTime)
		fmt.Printf("  avg_queue_time: %s\n", avgQueueTime)
		fmt.Printf("  max_queue_time: %s\n", station.maxQueueTime)
	}

	// Print statistics for registers.
	fmt.Println("Registers:")
	avgQueueTime := time.Duration(0)
	if registers.totalCars > 0 {
		avgQueueTime = registers.totalTime / time.Duration(registers.totalCars)
	}
	fmt.Printf("  total_cars: %d\n", registers.totalCars)
	fmt.Printf("  total_time: %s\n", registers.totalTime)
	fmt.Printf("  avg_queue_time: %s\n", avgQueueTime)
	fmt.Printf("  max_queue_time: %s\n", registers.maxQueueTime)
}

// chooseRandomStation selects a random station from available options.
func chooseRandomStation() string {
	stations := []string{"gas", "diesel", "lpg", "electric"}
	return stations[rand.Intn(len(stations))]
}
