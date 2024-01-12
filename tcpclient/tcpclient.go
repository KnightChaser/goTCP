package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Write user message on the data stream connected to the server
func SendMessageToServer(connection net.Conn, message string) {
	fmt.Fprintf(connection, "\n%s", message)
}

// Receive and print the message from the server
func ServerMessageReceiver(connection net.Conn, currentUsername string) {
	serverMessageHandler := bufio.NewScanner(connection)
	for serverMessageHandler.Scan() {
		message := serverMessageHandler.Text()
		fmt.Printf("%s\n", message)
	}
}

// argument: username(string)
func main() {

	username := flag.String("username", "exampleUser", "Your own username for chatting")
	flag.Parse()

	if *username == "" {
		fmt.Println("Please enter your username via -username option")
		return
	}

	protocol := "tcp"
	accessingAddressPort := "127.0.0.1:7777" // localhost
	connection, err := net.Dial(protocol, accessingAddressPort)

	if err != nil {
		fmt.Printf("Error while connecting %s(protocol: %s)\n", accessingAddressPort, protocol)
		log.Fatal(err)
		return
	} else {
		fmt.Printf("Connected to %s(protocol: %s)!\n", accessingAddressPort, protocol)
		fmt.Fprintf(connection, "%s\n", *username) // send the server about username
	}
	defer connection.Close()

	go ServerMessageReceiver(connection, *username)

	var userInput string

	// Scanner to capture user input
	userInputScanner := bufio.NewScanner(bufio.NewReader(os.Stdin))

	for {
		// If username is set, ready to chat!
		fmt.Printf("%s(You): ", *username)

		// Check if there's any input
		if userInputScanner.Scan() {
			userInput = userInputScanner.Text()

			// A special triggering keyword to terminate the user connection
			if strings.ToLower(userInput) == "!exit" {
				fmt.Printf("Terminating %s connection with the server(%s)\n",
					strings.ToUpper(protocol), accessingAddressPort)
				return
			}

			SendMessageToServer(connection, userInput)
		}
	}
}
