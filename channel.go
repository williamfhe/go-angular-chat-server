package main

import "log"

type Channel struct {
	hub     *Hub
	Name    string
	clients map[*Client]bool
}

func NewChannel(hub *Hub, name string) *Channel {
	log.Printf("New Channel: %s", name)
	return &Channel{
		hub:     hub,
		Name:    name,
		clients: make(map[*Client]bool),
	}
}

func (channel *Channel) AddClient(client *Client) {
	log.Printf("Channel: %s; new Client!\n", channel.Name)
	channel.clients[client] = true
}

func (channel *Channel) DeleteClient(client *Client) {
	delete(channel.clients, client)
	if len(channel.clients) == 0 {
		channel.hub.terminateChannel(channel.Name)
	}
}

func (channel *Channel) BroadcastBytes(bytes []byte) {
	for client := range channel.clients {
		select {
		case client.send <- bytes:
		default:
			channel.hub.deleteClient(client)
		}
	}
}

func (channel *Channel) Broadcast(notification *Notification) {
	for client := range channel.clients {
		if client.Id != notification.Client.Id {
			select {
			case client.send <- notification.Content:
			default:
				channel.hub.deleteClient(client)
			}
		}
	}
}
