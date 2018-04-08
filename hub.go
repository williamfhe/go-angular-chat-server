package main

import (
	"fmt"
	"log"
	"encoding/json"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients  map[*Client]bool
	channels map[string]*Channel

	// Inbound messages from the clients.
	broadcast chan Notification

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}


func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Notification),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channels:   make(map[string]*Channel),
		clients:    make(map[*Client]bool),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				hub.deleteClient(client)
			}
		case notification := <-hub.broadcast:
			fmt.Println(string(notification.Content), len(hub.clients))
			hub.dealWithNotification(&notification)
		}
	}
}

func (hub *Hub) addToChannel(chanName string, client *Client) {
	channel, ok := hub.channels[chanName]
	if !ok {
		hub.newChannel(chanName)
	}
	channel.AddClient(client)
}

func (hub *Hub) newChannel(name string) *Channel {
	channel := NewChannel(hub, name)
	hub.channels[name] = channel
	return channel
}

func (hub *Hub) terminateChannel(name string) {
	delete(hub.channels, name)
}

func (hub *Hub) deleteClient(client *Client) {
	var leaveNotification LeaveNotification

	leaveNotification.Type = "leaveChannel"
	leaveNotification.Data.User = client.Username

	for channelName := range client.Channels {
		if channel, ok := hub.channels[channelName]; ok {
			channel.DeleteClient(client)
			leaveNotification.Data.ChannelName = channelName

			bytes, err := json.Marshal(channelName)
			if err != nil {
				panic(err)
			}
			channel.BroadcastBytes(bytes)
		}
	}
	close(client.send)
	delete(hub.clients, client)
}

func (hub *Hub) broadcastOnChannel(channelName string, notification *Notification) {
	if channel, ok := hub.channels[channelName]; ok {
		channel.Broadcast(notification)
	} else {
		log.Printf("Channel: %s doesn't exist\n", channelName)
	}

}
