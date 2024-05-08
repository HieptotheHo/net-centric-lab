package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var clientList = make(map[*net.UDPAddr]string)

func printClientList() {
	fmt.Println("Client List:")
	for key, value := range clientList {
		fmt.Println(key, ":", value)
	}
}

func main() {
	// Create a UDP address
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	checkError(err)

	// Create a UDP listener
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	defer conn.Close()

	fmt.Println("Server is listening on port 8080...")

	// Handle incoming messages
	for {
		// Buffer to hold incoming data
		buf := make([]byte, 1024)
		// Read data from the connection
		n, addr, err := conn.ReadFromUDP(buf[0:])
		checkError(err)

		received := string(buf[:n])
		pieces := strings.Split(received, " ")
		fmt.Println(received)
		clientName := pieces[0]
		prefix := pieces[1]
		var message string
		if len(pieces) >= 3 {
			message = strings.Join(pieces[2:], " ")
		}

		if received == "exit" {
			delete(clientList, addr)
		} else {
			if prefix == "@name" {
				fmt.Println(clientName, "-", addr, "has been connected!")
				clientList[addr] = clientName
			} else {
				if prefix == "@all" {
					for address, name := range clientList {
						if address != addr && name != clientName {
							conn.WriteToUDP([]byte(clientName+": "+message), address)
						}
					}
				} else {
					receiver := prefix[1:]
					for address, name := range clientList {
						if receiver == name {
							conn.WriteToUDP([]byte(clientName+": "+message), address)
						}
					}
				}
			}

		}

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
