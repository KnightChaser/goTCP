package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/gosuri/uilive"
	"github.com/inancgumus/screen"
)

// Return the first 6 characters of the given string input
// A hashcode identifier of the users, might be an identifier of users with same name
func Sha256First6(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	result := hashString[:6]

	return result
}

func GenerateRandomNumber() int {
	randomNumber := rand.Intn(100) + 1

	return randomNumber
}

// Write user message on the data stream connected to the server
func SendMessageToServer(connection net.Conn, message string) {
	fmt.Fprintf(connection, "%s\n", message)
}

// Receive and print the message from the server,
// and print the data on the console appropriately
func ConsoleMessageHandler(connection net.Conn, currentUsername string) {
	var messageFromServerList []string
	serverMessageHandler := bufio.NewScanner(connection)
	consoleUILiveWriter := uilive.New() // auto refreshing feature the console
	consoleUILiveWriter.Start()

	for serverMessageHandler.Scan() {
		message := serverMessageHandler.Text()
		messageFromServerList = append(messageFromServerList, message)

		// Refreshing the console
		var currentlyPrintedDataOnConsole string
		for _, messageLine := range messageFromServerList {
			currentlyPrintedDataOnConsole = fmt.Sprintf("%s\n%s", currentlyPrintedDataOnConsole, messageLine)
		}
		// Provide user interface to continue typing...
		currentlyPrintedDataOnConsole = fmt.Sprintf("%s\nYou(%s)> ", currentlyPrintedDataOnConsole, currentUsername)
		fmt.Fprintf(consoleUILiveWriter, currentlyPrintedDataOnConsole)
		consoleUILiveWriter.Flush()
	}
	consoleUILiveWriter.Stop()
}

// argument: username(string)
func main() {

	username := flag.String("username", "exampleUser", "Your own username for chatting")
	flag.Parse()

	if *username == "" {
		fmt.Println("Please enter your username via -username option")
		return
	}

	// Add a random seed value
	randomSeed := strconv.Itoa(GenerateRandomNumber())
	*username = fmt.Sprintf("%s#%s", *username, Sha256First6(randomSeed))

	protocol := "tcp"
	accessingAddressPort := "127.0.0.1:7777" // localhost
	// accessingAddressPort := "192.168.111.111:7777" // docker standard
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

	go ConsoleMessageHandler(connection, *username)

	var userInput string

	// Scanner to capture user input
	userInputScanner := bufio.NewScanner(os.Stdin)

	// If username is set, ready to chat!
	fmt.Printf("You are \"%s\". Hit the first message! > ", *username)

	for {

		// Check if there's any input
		if userInputScanner.Scan() {
			userInput = userInputScanner.Text()

			// A special triggering keyword to terminate the user connection
			if strings.ToLower(userInput) == "!exit" {
				fmt.Printf("Terminating %s connection with the server(%s)\n",
					strings.ToUpper(protocol), accessingAddressPort)
				return
			}

			screen.Clear()
			screen.MoveTopLeft()
			SendMessageToServer(connection, userInput)
		}
	}
}
