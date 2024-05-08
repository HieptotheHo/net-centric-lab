package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

type Credentials struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type UserList struct {
// 	Users []User `json:"users"`
// }

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func readFromJSONFile(filename string) ([]Credentials, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var credentials []Credentials
	err = json.NewDecoder(file).Decode(&credentials)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return credentials, nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func authenticate(credentials []Credentials, clientCredentials Credentials) Response {
	for _, cred := range credentials {
		if cred.Email == clientCredentials.Email {
			// Compare the hashed password with the stored hashed password
			err := bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(clientCredentials.Password))
			if err == nil {
				return Response{Success: true, Message: "Authentication successful"}
			}
		}
	}
	return Response{Success: false, Message: "Invalid credentials"}
}

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

func chooseWord() string {
	words := []string{"apple", "banana", "orange", "grape", "kiwi", "strawberry"}

	return words[rand.Intn(len(words))]
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
func displayWord(word string, guessedLetters []string) string {
	displayedWord := ""
	for _, letter := range word {
		if contains(guessedLetters, string(letter)) {
			displayedWord += string(letter)
		} else {
			displayedWord += "_"
		}
	}
	return displayedWord
}
func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	credentials, err := readFromJSONFile("credentials.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var clientCredentials Credentials
	// Print the content of the buffer for debugging
	fmt.Println("Buffer:", string(buffer))

	// Trim null characters from the buffer
	trimmedBuffer := bytes.TrimRight(buffer, "\x00")

	// Decode JSON from trimmed buffer
	err = json.Unmarshal(trimmedBuffer, &clientCredentials)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	response := authenticate(credentials, clientCredentials)
	fmt.Println(response)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		return
	}

	_, err = conn.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error sending response to client:", err)
		return
	}
	if response.Success {

		word := chooseWord()
		fmt.Println("Secret word: ", word)
		var guessedLetters []string
		// attempts := 6

		///////////////////
		// write data to response
		time := time.Now().Format(time.ANSIC)
		session := rand.Intn(999)
		responseStr := fmt.Sprintf("\nSuccessful Connection at %v\n\nWelcome to Hangman!\nTry to guess the word\n%s", time, displayWord(word, guessedLetters))

		_, err = fmt.Fprintf(conn, "\n"+responseStr)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		_, err = fmt.Fprintf(conn, "%d", session)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		attempts := 6
		for attempts > 0 {

			reader := bufio.NewReader(conn)
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from client:", err)
				return
			}
			fmt.Print("Client: ", message)
			// Extract session and guessNum
			parts := strings.Split(strings.TrimSpace(message), "_")
			guess := strings.TrimSpace(parts[1])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(guess)

			if len(guess) != 1 || !strings.ContainsAny(guess, "abcdefghijklmnopqrstuvwxyz") {
				fmt.Fprintln(conn, "Please enter a single letter.")
				continue
			}

			if contains(guessedLetters, guess) {
				fmt.Fprintln(conn, "You've already guessed that letter.")
				continue
			}

			guessedLetters = append(guessedLetters, guess)

			fmt.Println(guessedLetters)
			if !strings.Contains(word, guess) {
				attempts--
				fmt.Fprintf(conn, "Wrong guess! You have %d attempts left.\n", attempts)
				if attempts == 0 {
					fmt.Fprintf(conn, "You're out of attempts! The word was: %s\n", word)
					break
				}
			} else {
				fmt.Fprintln(conn, "Correct guess!")
			}
			if strings.Contains(displayWord(word, guessedLetters), "_") == false {
				fmt.Fprintf(conn, "Congratulations! You guessed the word: %s\n", word)
				break
			}
		}

	}

}
