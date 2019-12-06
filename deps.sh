#!/bin/bash
# Author:       nghiatc
# Email:        congnghia0609@gmail.com

#source /etc/profile

echo "Install library dependencies..."
go get -u github.com/tools/godep
go get -u github.com/nats-io/nats-server
go get -u github.com/nats-io/nats.go/

echo "Install dependencies complete..."
