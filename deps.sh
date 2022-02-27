#!/bin/bash
# Author:       nghiatc
# Email:        congnghia0609@gmail.com

#source /etc/profile

echo "Install library dependencies..."

go get -u github.com/tools/godep
go get -u github.com/congnghia0609/ntc-gconf
#go get -u github.com/nats-io/nats-server
go get -u github.com/nats-io/nats-server/v2
go get -u github.com/nats-io/nats.go/
go get -u github.com/shettyh/threadpool
go get -u github.com/google/uuid

echo "Install dependencies complete..."
