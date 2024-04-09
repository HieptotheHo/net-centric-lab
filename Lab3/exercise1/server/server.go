package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(100) + 1
	fmt.Println("The secret number:", secretNumber)
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// write data to response
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("Successful Connection at %v", time)
	conn.Write([]byte(responseStr))

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		data := buffer[:n]
		fmt.Println("Received:", string(data))
		str := string(data)

		// Parse string to integer
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		var returnMessage string
		if num > secretNumber {
			returnMessage = "Your number is greater than mine."
		} else if num < secretNumber {
			returnMessage = "Your number is smaller than mine."
		} else {
			returnMessage = "You guess is correct!!!\n------\nServer has changed its secret number. Take new guess!"
			secretNumber = rand.Intn(100) + 1
			fmt.Println("The secret number has been changed to:", secretNumber)
		}
		_, err = conn.Write([]byte(returnMessage))
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}

}
