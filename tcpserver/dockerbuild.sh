#!/bin/bash

sudo docker build . --tag gotcpserver:0.1
sudo docker run --name gotcpserver -p 7777:7777 gotcpserver:0.1