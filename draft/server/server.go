package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var clientList = make(map[*net.UDPAddr]string)

func printClientList() {
	fmt.Println("#### CLIENT LIST ####")
	for key, value := range clientList {
		fmt.Println(key, ":", value)
	}
	fmt.Println("#####################")
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

				dup := false
				for _, name_temp := range clientList {
					if clientName == name_temp {
						conn.WriteToUDP([]byte("Your name has been registered before!"), addr)
						dup = true
					}
				}

				fmt.Println("DEBUG: ", dup)
				if dup == false {
					fmt.Println(clientName, "-", addr, "has been connected!")
					clientList[addr] = clientName
					conn.WriteToUDP([]byte("Registered successfully!"), addr)
					printClientList()

				}

			} else if prefix == "@all" {
				for address, name := range clientList {
					if address != addr && name != clientName {
						conn.WriteToUDP([]byte(clientName+": "+message), address)
					}
				}
			} else if strings.HasPrefix(prefix, "@") {
				receiver := prefix[1:]
				available := false
				for address, name := range clientList {
					if receiver == name {
						conn.WriteToUDP([]byte(clientName+": "+message), address)
						available = true
						break
					}
				}
				if available == false {
					conn.WriteToUDP([]byte("No client name matched!"), addr)
				}
			} else {
				conn.WriteToUDP([]byte("Server: Wrong Syntax!"), addr)
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
