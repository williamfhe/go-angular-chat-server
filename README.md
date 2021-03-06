# GoNg Chat

The server from the go-angular websocket chat.

The client can be found [here](https://github.com/williamfhe/go-angular-chat-client).

## Functionalities

* Choose a Username
* Join a Channel
* Talk with other users
* Leave the channel
* Get notifications when people join or leave the channel

## Installation

You are gonna need Go to launch the server, if it isn't already installed click [here](https://golang.org/).

Clone this repository

```bash
go get github.com/williamfhe/go-angular-chat-server
cd $GOPATH/src/github.com/williamfhe/go-angular-chat-server
```

Build the server and launch it

```bash
go build .
./go-angular-chat-server
```

## Docker

```bash
docker build -t gong-server .
docker run -p 8080:8080 gong-server
```

## Usage

Open your navigator and go to http://localhost:8080/ or click [here](http://localhost:8080/).
