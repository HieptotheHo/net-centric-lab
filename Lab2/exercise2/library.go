package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	libraryCapacity = 30
	totalStudents   = 100
	minVisitTime    = 1
	maxVisitTime    = 4
)

var (
	wg           sync.WaitGroup
	library      chan struct{}
	waitingMutex sync.Mutex
	waitingList  []int
)

func main() {
	rand.Seed(time.Now().UnixNano())

	library = make(chan struct{}, libraryCapacity)

	for i := 0; i < totalStudents; i++ {
		wg.Add(1)
		go student(i)
	}

	wg.Wait()
	close(library)
	fmt.Printf("The library was open for %d hours to serve all the students.\n", totalStudents*maxVisitTime)
}

func student(id int) {
	defer wg.Done()

	// Simulate student's visit duration
	visitTime := rand.Intn(maxVisitTime-minVisitTime+1) + minVisitTime

	// Enter library (wait if library is full)
	select {
	case library <- struct{}{}:
		fmt.Printf("Time 0: Student %d entered the library.\n", id)
	default:
		// Student has to wait
		waitingMutex.Lock()
		waitingList = append(waitingList, id)
		fmt.Printf("Time 0: Student %d is waiting.\n", id)
		waitingMutex.Unlock()
		return // Student can't proceed further until they enter the library
	}

	// Simulate student reading in the library
	for i := 1; i <= visitTime; i++ {
		time.Sleep(time.Hour)
		fmt.Printf("Time %d: Student %d is reading at the library.\n", i, id)
	}

	// Leave library
	<-library
	fmt.Printf("Time %d: Student %d left the library after spending %d hours.\n", visitTime, id, visitTime)

	// Check if there are students waiting
	waitingMutex.Lock()
	if len(waitingList) > 0 {
		numToEnter := rand.Intn(min(len(waitingList), libraryCapacity-len(library))) + 1 // Random number of students to enter
		for j := 0; j < numToEnter; j++ {
			studentID := waitingList[0]
			waitingList = waitingList[1:]
			library <- struct{}{}
			fmt.Printf("Time %d: Student %d started reading after waiting.\n", visitTime, studentID)
		}
	}
	waitingMutex.Unlock()
}
