package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	//READ INPUT AS PORT
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}

	//CREATE UDP ADDRESS
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)
	fmt.Println("hello")
	//CONNECT TO SERVER USING UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	//GREETING SERVER
	_, err = conn.Write([]byte("Hello, Server!"))
	checkError(err)

	//READ FROM SERVER
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)
	fmt.Println("hello2")
	fmt.Println(string(buf[0:n]))

	//SEND NAME TO SERVER
	var name string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)
	//send
	_, err = conn.Write([]byte(name))
	checkError(err)

	var message string
	for {
		fmt.Print("Enter your message: ")
		fmt.Scanln(&message)
		_, err = conn.Write([]byte(message))
		checkError(err)

		n, err := conn.Read(buf[0:])
		checkError(err)

		fmt.Println(string(buf[0:n]))
	}

	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
