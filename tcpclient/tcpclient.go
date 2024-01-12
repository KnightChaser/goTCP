package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Write user message on the data stream connected to the server
func sendMessageToServer(connection net.Conn, message string) {
	fmt.Fprintf(connection, "%s\n", message)
}

// Receive and print the server message
func serverMessageReceiver(connection net.Conn, serverMessageChannel chan string) {
	serverMessageHandler := bufio.NewScanner(connection)
	for serverMessageHandler.Scan() {
		message := serverMessageHandler.Text()
		serverMessageChannel <- message
	}
	close(serverMessageChannel)
}

func main() {
	protocol := "tcp"
	accessingAddressPort := "127.0.0.1:7777" // localhost
	connection, err := net.Dial(protocol, accessingAddressPort)
	if err != nil {
		fmt.Printf("Error while connecting %s(protocol: %s)\n", accessingAddressPort, protocol)
		log.Fatal(err)
		return
	}
	defer connection.Close()

	// Read messages from the server, using goroutine
	serverMessageChannel := make(chan string)
	go serverMessageReceiver(connection, serverMessageChannel)

	// Print server messages completely first (via channel, synchronously)
	go func() {
		for message := range serverMessageChannel {
			fmt.Printf("Server says: %s\n", message)
		}
	}()

	// Read the user input and send messages to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		// A special triggering keyword to terminate user connection
		if strings.ToLower(message) == "!exit" {
			fmt.Printf("Terminating %s connection with the server(%s)\n",
				strings.ToUpper(protocol), accessingAddressPort)
			break
		}

		sendMessageToServer(connection, message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading user input")
		log.Fatal(err)
		return
	}
}
