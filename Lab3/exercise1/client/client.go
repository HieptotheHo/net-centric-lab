package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("Hello, Server!"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(received), "\n")
	//////////////////////////////////////////////////////////////
	fmt.Println("Welcome to Guessing Game!")
	fmt.Println("Now start guessing a number in the range of 1 to 100\n")
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a number: ")
		scanner.Scan()
		text := scanner.Text()
		if text == "exit" {
			break
		}
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			return
		}
		fmt.Println("Server response:", string(response[:n]))
	}
}
