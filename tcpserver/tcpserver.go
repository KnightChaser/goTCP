package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Handle client messages toward this server
func clientHandler(connection net.Conn, clientRemoteAddress net.Addr) {
	defer connection.Close()

	clientDataReceiver := bufio.NewScanner(connection)
	for clientDataReceiver.Scan() {
		message := clientDataReceiver.Text()
		fmt.Printf("[~] Client(%v) said => %s\n", clientRemoteAddress, message)
	}

	if err := clientDataReceiver.Err(); err != nil {
		fmt.Printf("[!] Client(%v) disconnected abnormally, error: %v\n", clientRemoteAddress, err)
	} else {
		fmt.Printf("[-] Client(%v) disconnected normally\n", clientRemoteAddress)
	}
}

// Send a welcome message to the server
func sendWelcomeMessageToClients(connection net.Conn) {
	fmt.Fprintf(connection, "Hello, user!, welcome to the server.\n")
}

func main() {
	protocol := "tcp"
	listeningAddressPort := "127.0.0.1:7777" // localhost
	listener, err := net.Listen(protocol, listeningAddressPort)
	if err != nil {
		fmt.Println("[!] Error while listening the server")
		log.Fatal(err)
		return
	}
	defer listener.Close()

	fmt.Printf("[O] The server is now listening from %s(protocol: %s)\n", listeningAddressPort, protocol)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("[!] Err while accepting the client request(s)")
			log.Fatal(err)
			break
		}
		fmt.Printf("[+] A new connection from: %v\n", connection.RemoteAddr())
		sendWelcomeMessageToClients(connection)
		go clientHandler(connection, connection.RemoteAddr()) // Handle user request as goroutine
	}

	fmt.Printf("[X] Server service(%s(protocol: %s)) terminated.\n", listeningAddressPort, protocol)
}
