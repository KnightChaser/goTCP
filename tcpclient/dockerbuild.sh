#!/bin/bash

if [[ $(/usr/bin/id -u) -ne 0 ]]; then
        echo "Not running as root"
        exit
fi

# Get the USERNAME from user input
read -p "Enter USERNAME: " USERNAME
echo "You're now @$USERNAME"

# Generate a 6-digit random hex string
RandomHex=$(cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 6 | head -n 1)
containerName="gotcpclient${USERNAME}${RandomHex}"

# Build and run Docker container
docker build -t gotcpclient:0.1 --build-arg --no-cache .
docker run --rm --network goTCPnet --interactive --tty --name "$containerName" gotcpclient:0.1 ./tcpclient -username "$USERNAME"

