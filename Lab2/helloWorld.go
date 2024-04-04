package main

import (
	"fmt"
	"math/rand"
	"time"
)

func BufferedChannel() {
	// Create a buffered channel with a capacity of 2
	cars := make(chan string, 2)

	// Start two Goroutines that add cars to the channel
	go addCar("Ferrari", cars)
	go addCar("Lamborghini", cars)

	// Wait for the cars to be added to the channel
	time.Sleep(2 * time.Second)
	close(cars)

	// Start a Goroutine that simulates the race
	go startRace(cars)

	// Wait for the race to finish
	time.Sleep(6 * time.Second)

	fmt.Println("Race over!")
}

func addCar(name string, cars chan string) {
	cars <- name
	fmt.Println(name, "added to the race!")
}

func startRace(cars chan string) {
	for {
		// Receive a car from the channel
		car, open := <-cars
		if !open {
			break
		}

		fmt.Println(car, "is racing...")

		// Simulate the race by waiting for a random duration
		time.Sleep(time.Duration(1+rand.Intn(5)) * time.Second)
	}

	fmt.Println("All cars have finished the race!")
}
