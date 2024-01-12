package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

// An object representing a single connected user
type User struct {
	connection net.Conn
	username   string
}

// A connect client pool for management
// This object might be shared by multiple goroutines,
// so consistency protection by MUTEX is highly encouraged.
type ConnectedClientPool struct {
	clients    map[net.Addr]User
	clientsQty uint32
	mutex      sync.Mutex
}

// Initially, createa a client pool and set clientsQty as zero(0; "no one joined yet").
func CreateConnectedClientPool() *ConnectedClientPool {
	return &ConnectedClientPool{
		clients:    make(map[net.Addr]User),
		clientsQty: 0,
	}
}

// Adds a new connected user to the pool
func (connectedClientPool *ConnectedClientPool) AddUserToConnectedClientPool(clientRemoteAddress net.Addr, user User) {
	connectedClientPool.mutex.Lock()
	defer connectedClientPool.mutex.Unlock()

	connectedClientPool.clients[clientRemoteAddress] = user
	connectedClientPool.clientsQty += 1
}

// Deletes a new disconnected user from the pool
func (connectedClientPool *ConnectedClientPool) DeleteUserFromConnectedClientPool(clientRemoteAddress net.Addr) {
	connectedClientPool.mutex.Lock()
	defer connectedClientPool.mutex.Unlock()

	delete(connectedClientPool.clients, clientRemoteAddress)
	connectedClientPool.clientsQty -= 1
}

// Gets an user object by username
func (connectedClientPool *ConnectedClientPool) GetUserByUsername(username string) (User, bool) {
	connectedClientPool.mutex.Lock()
	defer connectedClientPool.mutex.Unlock()

	// Just search one by one, O(N)
	for _, user := range connectedClientPool.clients {
		if username == user.username {
			return user, true
		}
	}
	return User{}, false
}

// Broadcast received message to every user except for the OP (such like chatroom)
func (connectedClientPool *ConnectedClientPool) BroadcastMessage(senderAddress net.Addr, message string) {
	connectedClientPool.mutex.Lock()
	defer connectedClientPool.mutex.Unlock()

	for userRemoteAddress, user := range connectedClientPool.clients {
		// Exclude OP
		if senderAddress != userRemoteAddress {
			fmt.Fprintf(user.connection, "%s\n", message)
		}
	}
}

// Handle client messages toward this server
func ClientHandler(connection net.Conn, clientRemoteAddress net.Addr, connectedClientPool *ConnectedClientPool) {
	// When the user terminates its session, proceed with the cleaning user information procedure
	// whether the user disconnected its connection normally or not
	defer func() {
		connection.Close()
		connectedClientPool.DeleteUserFromConnectedClientPool(clientRemoteAddress)
	}()

	// Ask username to the newly joined user and register user information to connectedClientPool
	// The first input of tcpclient.go will be the username (guided on the tcpclient.go)
	userNameScanner := bufio.NewScanner(connection)
	userNameScanner.Scan()
	username := userNameScanner.Text()
	user := User{
		connection: connection,
		username:   username,
	}
	connectedClientPool.AddUserToConnectedClientPool(clientRemoteAddress, user)

	// Alerting everyone for a new user join
	fmt.Printf("[+] User \"%s\"(from %v) just joined the server.\n", username, clientRemoteAddress)
	connectedClientPool.BroadcastMessage(clientRemoteAddress, fmt.Sprintf("User \"%s\" just joined the chat!", username))

	// Registration finished, accepting the user to interact
	clientDataReceiver := bufio.NewScanner(connection)
	for clientDataReceiver.Scan() {
		message := clientDataReceiver.Text()
		fmt.Printf("[~] Client(@%s/%v) said => %s\n", username, clientRemoteAddress, message)
		connectedClientPool.BroadcastMessage(clientRemoteAddress, fmt.Sprintf("%s: %s", username, message))
	}

	if err := clientDataReceiver.Err(); err != nil {
		fmt.Printf("[!] Client(@%s/%v) disconnected abnormally, error: %v\n", username, clientRemoteAddress, err)
		connectedClientPool.BroadcastMessage(clientRemoteAddress, fmt.Sprintf("User \"%s\" just left the chat (normally)", username))
	} else {
		fmt.Printf("[-] Client(@%s/%v) disconnected normally\n", username, clientRemoteAddress)
		connectedClientPool.BroadcastMessage(clientRemoteAddress, fmt.Sprintf("User \"%s\" just left the chat (normally)", username))
	}
}

// Send specific message to the client(specific client, not broadcasting)
func SendMessageToClients(connection net.Conn, message string) {
	fmt.Fprintf(connection, "%s\n", message)
}

func main() {
	// Open up the server
	protocol := "tcp"
	// listeningAddressPort := "127.0.0.1:7777" // localhost
	listeningAddressPort := "192.168.111.111:7777"
	listener, err := net.Listen(protocol, listeningAddressPort)
	if err != nil {
		fmt.Println("[!] Error while listening the server")
		log.Fatal(err)
		return
	}
	defer listener.Close()

	// Ready to open the server
	connectedClientPool := CreateConnectedClientPool()
	fmt.Printf("[O] The server is now listening from %s(protocol: %s)\n", listeningAddressPort, protocol)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("[!] Err while accepting the client request(s)")
			log.Fatal(err)
			break
		}

		clientRemoteAddress := connection.RemoteAddr()
		fmt.Printf("[+] A new connection from: %v\n", clientRemoteAddress)

		go ClientHandler(connection, clientRemoteAddress, connectedClientPool) // Handle user request as goroutine
	}

	fmt.Printf("[X] Server service(%s(protocol: %s)) terminated.\n", listeningAddressPort, protocol)
}
