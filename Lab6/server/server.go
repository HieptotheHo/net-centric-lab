package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Listen on port 9999
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started on port 9999")

	for {
		// Accept connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle request in a separate goroutine
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Read the request
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// Parse the request line
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		fmt.Println("Malformed request line:", requestLine)
		return
	}
	method, url, version := parts[0], parts[1], parts[2]

	// Only handle GET requests
	if method != "GET" {
		fmt.Println("Unsupported method:", method)
		return
	}

	// Remove leading slash from URL to get the file path
	if url == "/" {
		url = "/index.html"
	}
	filePath := "." + url

	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		sendResponse(conn, version, 404, "Not Found", "text/plain", []byte("404 Not Found"))
		return
	}
	defer file.Close()

	// Read the file content
	fileInfo, err := file.Stat()
	if err != nil {
		sendResponse(conn, version, 500, "Internal Server Error", "text/plain", []byte("500 Internal Server Error"))
		return
	}

	content := make([]byte, fileInfo.Size())
	_, err = file.Read(content)
	if err != nil {
		sendResponse(conn, version, 500, "Internal Server Error", "text/plain", []byte("500 Internal Server Error"))
		return
	}

	// Send the response
	sendResponse(conn, version, 200, "OK", "text/html", content)
}

func sendResponse(conn net.Conn, version string, statusCode int, statusText, contentType string, body []byte) {
	responseLine := fmt.Sprintf("%s %d %s\r\n", version, statusCode, statusText)
	headers := fmt.Sprintf("Content-Length: %d\r\nContent-Type: %s\r\n\r\n", len(body), contentType)
	response := responseLine + headers
	conn.Write([]byte(response))
	conn.Write(body)
}
