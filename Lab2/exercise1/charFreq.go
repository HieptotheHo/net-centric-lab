package main

import (
	"fmt"
	"sync"
)

func countCharacters(input string, wg *sync.WaitGroup, ch chan<- map[rune]int) {
	defer wg.Done()

	charCount := make(map[rune]int)
	for _, char := range input {
		charCount[char]++
	}

	ch <- charCount
}

func main() {
	input := "There is no one who loves pain itself, who seeks after it and wants to have it, simply because it is pain..."
	numberOfChunks := 12

	// Create channels and wait group
	ch := make(chan map[rune]int) // buffered channel to prevent deadlock
	var wg sync.WaitGroup

	// Split the string for concurrent processing
	chunkSize := len(input) / numberOfChunks // Splitting the string into 4 chunks
	var startIndex, endIndex int
	for i := 0; i < numberOfChunks; i++ {
		startIndex = i * chunkSize
		endIndex = (i + 1) * chunkSize
		if i == numberOfChunks-1 {
			endIndex = len(input)
		}
		wg.Add(1)
		go countCharacters(input[startIndex:endIndex], &wg, ch)
	}

	// Collect and merge the results
	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make(map[string]int)
	for charCount := range ch {
		for char, count := range charCount {
			// Convert rune to string
			result[string(char)] += count
		}
	}

	// Print the result of character counts
	fmt.Println("Character counts:", result)
}
