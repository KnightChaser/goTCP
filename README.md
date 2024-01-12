# goTCP
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![PowerShell](https://img.shields.io/badge/PowerShell-%235391FE.svg?style=for-the-badge&logo=powershell&logoColor=white)
![Shell Script](https://img.shields.io/badge/shell_script-%23121011.svg?style=for-the-badge&logo=gnu-bash&logoColor=white)

A simple Go chatting application with TCP socket, based on the client-server model

> Demonstration(picture)
![img](./READMEasset/figure1.png)

This simple project is a simple real-time chat application(because it's simple, please consider the naive interfaces.) implemented in the Go programming language. It consists of a server (`tcpserver/tcpserver.go`) and a client (`tcpclient/tcpclient.go`) that communicate over a TCP network. The application allows users to connect to the server, set a username, and exchange messages in real time.

### Feature
* **Server-Client Architecture**: The application follows a server-client architecture, enabling multiple users to connect to the server simultaneously.

* **Real-Time Messaging**: Users can exchange messages in real time. When a user sends a message, it is immediately broadcast to all other connected users.

* **Username Registration**: Users are prompted to set a username when connecting to the server. This username is used to identify users in the chat.

* **Graceful Termination**: Users can terminate their connection gracefully by typing `!exit` in the client interface, triggering a clean disconnection from the server.

* **Dockerized**: Server and client are encapsulated by docker. Because there are few lines of docker-related commands for every instance, this project provides a "one-stop script" that processes setting up, docker building, and docker execution at once
    * [x] Windows(Powershell script(`dockerbuild.ps1`))
    * [x] Linux(Shell script(`dockerbuild.sh`))


### Execution
Because this project's theme is communicating with each other, the docker networking setup is very important. You just need to execute the "one-stop script" that this project provides. **Support both Windows and Linux!**
* **Server**: `tcpserver\dockerbuild.ps1` if **Windows**, `sudo bash tcpserver/dockerbuild.sh` if **Linux** (Only one instance is required, so execute this script only once at the same time.)
* **Client**: `tcpclient\dockerbuild.ps1` if **Windows**, `sudo bash tcpclient/dockerbuild.sh` if **Linux** (You may create client instances as many as you want.)

### Docker setup explanation (networking)
* **Server** (Windows Powershell)
```powershell
docker network remove goTCPnet
docker network create goTCPnet --subnet=192.168.111.0/24
docker build . --tag gotcpserver:0.1
docker run --rm --name gotcpserver --network goTCPnet --ip 192.168.111.111 -p 7777:7777 gotcpserver:0.1
```
Because the server accepts clients' requests remotely, IP address and port settings should be synchronized. By default, server setup creates a docker network `goTCPnet` and designates the network range as `192.168.111.0/24`, server process address as `192.168.111.111(/24)` manually. This setting should be synchronized with `tcpserver\tcpserver.go`'s `listeningAddressPort` variable that the server process refers to for setting up the address to accept clients' requests.
```go
listeningAddressPort := "192.168.111.111:7777"
```

* **Client** (Windows Powershell)
```powershell
# Get the USERNAME from user input
$USERNAME = Read-Host "Enter USERNAME "
Write-Output "You're now @$USERNAME"

# Generate a 6-digit random hex string (to make every container have different container names)
$RandomHex = -join ((1..6) | ForEach-Object { Get-Random -Minimum 0 -Maximum 16 } | ForEach-Object { $_.ToString("X") })
$containerName = "gotcpclient${USERNAME}${RandomHex}"

# Build and run Docker container
docker build -t gotcpclient:0.1 --build-arg --no-cache .
docker run --rm --network goTCPnet --interactive --tty --name $containerName gotcpclient:0.1 ./tcpclient -username $USERNAME
```
Because the client (`tcpclient\tcpclient.go`) accepts username as an argument `-username`, It asks a user for the username for every single created instance and applies them to the execution arguments. Suppose users may use multiple same names(`-username` argument). In that case, there might be some unexpected docker container name collisions, so client setup appends a 6-digit random hexadecimal string like `gotcpclient${USERNAME}${RandomHex}` to ensure every client container name should be globally unique. Client processes should be interactive to receive and print messages from users and server instances, `--interactive` and `--tty` options should be provided too.

### Note and contribution
* Because this was made for learning how to make "something" with Go, Docker, Socket communication, and Github, the implementation is naive and not-so-professional, I kindly ask for your consideration.
* Any contribution to this project is greatly appreciated.

<br>

**@KnightChaser**, a pro spaghetti chef in Gachon University, Cybersecurity department.