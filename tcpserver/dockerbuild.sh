#!/bin/bash

if [[ $(/usr/bin/id -u) -ne 0 ]]; then
	echo "Not running as root"
	exit
fi

docker network remove goTCPnet
docker network create goTCPnet --subnet=192.168.111.0/24
docker build . --tag gotcpserver:0.1
docker run --rm --name gotcpserver --network goTCPnet --ip 192.168.111.111 -p 7777:7777 gotcpserver:0.1
