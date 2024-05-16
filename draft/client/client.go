package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var clientName string

func main() {
	// Create a UDP address
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	checkError(err)

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	checkError(err)
	defer conn.Close()

	//receive message for register name
	buffer := make([]byte, 1024)
	for {
		fmt.Println("hello")
		var name string
		fmt.Print("Enter your name: ")
		fmt.Scanln(&name)
		clientName = name

		fmt.Println(name)
		//register name
		_, err = conn.Write([]byte(name + " " + "@name"))
		checkError(err)

		n, _, err := conn.ReadFromUDP(buffer)
		checkError(err)

		fmt.Println(string(buffer[:n]))
		if string(buffer[:n]) == "Registered successfully!" {
			break
		}
	}

	go receiveMessages(conn)
	sendMessages(conn)

}

func receiveMessages(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		checkError(err)

		fmt.Print("\n", string(buffer[:n]), "\n> ")
	}
}

func sendMessages(conn *net.UDPConn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		message := scanner.Text()
		_, err := conn.Write([]byte(clientName + " " + message))
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
