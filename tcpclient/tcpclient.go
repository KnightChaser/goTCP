package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Write user message on the data stream connected to the server
func SendMessageToServer(connection net.Conn, message string) {
	fmt.Fprintf(connection, "%s\n", message)
}

// Receive and print the message from the server
func ServerMessageReceiver(connection net.Conn, serverMessageChannel chan string) {
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
	} else {
		fmt.Printf("Connected to %s(protocol: %s)!\n", accessingAddressPort, protocol)
		fmt.Println("Type anything to start!")
	}
	defer connection.Close()

	// Read messages from the server, using goroutine
	serverMessageChannel := make(chan string)
	go ServerMessageReceiver(connection, serverMessageChannel)

	// Read the user input and send messages to the server
	// select-casing go channels
	for {
		select {
		case message, ok := <-serverMessageChannel:
			if !ok {
				// Server messages channel closed, exit
				fmt.Printf("The connection to the server(%s(protocol: %s)) was closed.\n", accessingAddressPort, protocol)
				return
			}
			fmt.Printf("Server says: %s\n", message)

		default:
			// Read user input and send messages to the server
			var userInput string
			fmt.Print("user input> ")
			fmt.Scanln(&userInput)

			// A special triggering keyword to terminate user connection
			if strings.ToLower(userInput) == "!exit" {
				fmt.Printf("Terminating %s connection with the server(%s)\n",
					strings.ToUpper(protocol), accessingAddressPort)
				return
			}

			SendMessageToServer(connection, userInput)
		}
	}
}
