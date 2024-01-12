#!/bin/bash

# Get the USERNAME from user input
read -p "Enter USERNAME: " USERNAME

sudo docker build --build-arg USERNAME="$USERNAME" -t gotcpclient:0.1 .
sudo docker run --rm --name gotcpclient gotcpclient:0.1
