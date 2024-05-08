package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {

		handleClient(conn)
	}

}

var clientList = make(map[*net.UDPAddr]string)

func printClientList() {
	fmt.Println("Client List:")
	for key, value := range clientList {
		fmt.Println(key, ":", value)
	}
}

func handleClient(conn *net.UDPConn) {
	//BUFFER DECLARATION
	var buf [512]byte

	//RESOLVE UDP ADDRESS
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	//SEND DATE TIME TO USER
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)

	//READNAME
	n, err := conn.Read(buf[0:])
	if err != nil {
		return
	}
	userName := string(buf[0:n])
	fmt.Println(addr, " - ", "\"", userName, "\"", " connected!")
	clientList[addr] = userName

	//WAIT FOR CLIENT'S MESSAGE
	for {
		n, err = conn.Read(buf[0:])
		if err != nil {
			return
		}
		message := string(buf[0:n])
		fmt.Println(message)

		for ipAddr, name := range clientList {
			message := name + ": " + message
			if name != userName {
				conn.WriteToUDP([]byte(message), ipAddr)
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
