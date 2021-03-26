#!/bin/bash
## for pre golang v1.16 like for Ubuntu 20.10 in my POC. For Ubuntu 21.04 using golang 1.16 in repo this is broken try 1.15 as below
sudo apt install git -y
sudo apt install golang -y
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
#sudo apt install golang-gopkg-mgo.v2-dev

## install golang 1.15 since 1.16 has issues with GOPATH being deprecated for modules.
#sudo wget https://golang.org/dl/go1.15.10.linux-amd64.tar.gz
#sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.15.10.linux-amd64.tar.gz

## mongodb admin tool
sudo snap install robo3t-snap

wget -qO - https://www.mongodb.org/static/pgp/server-4.4.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/4.4 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.4.list
sudo apt-get update
sudo apt-get install -y mongodb-org
ps --no-headers -o comm 1
sudo systemctl start mongod
sudo systemctl status mongod
sudo systemctl enable mongod

#go run Insert.go
#go run FindOne.go -id=<ObjectId>
#go run server.go -listen_address=192.168.1.225:8080
