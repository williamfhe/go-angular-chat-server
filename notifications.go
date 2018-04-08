package main

import (
	"github.com/tidwall/gjson"
	"log"
	"encoding/json"
)

type Notification struct {
	Content []byte
	Client  *Client
}

type ChannelInfoNotification struct {
	Type string `json:"type"`
	Data struct {
		ChannelName string `json:"channel"`
		Users []string `json:"users"`
	} `json:"data"`

}

type LeaveNotification struct {
	Type string `json:"type"`
	Data struct {
		ChannelName string `json:"channel"`
		User string `json:"user"`
	} `json:"data"`
}

func (hub *Hub) dealWithNotification(notification *Notification) {
	notificationType := gjson.GetBytes(notification.Content, "type").String()
	switch notificationType {
	case "newMessage":
		hub.dealWithUserMessage(notification)
	case "joinChannel":
		hub.userJoinChannel(notification)
	case "leaveChannel":
		hub.userLeaveChannel(notification)
	case "setUsername":
		hub.setClientUsername(notification)
	default:
		log.Print("Unrecognized json: ", string(notification.Content))

	}
}

func (hub *Hub) userLeaveChannel(notification *Notification) {
	channelName := gjson.GetBytes(notification.Content, "data.channel").String()
	channel := hub.channels[channelName]

	notification.Client.DeleteChannel(channelName)
	channel.DeleteClient(notification.Client)

	channel.Broadcast(notification)
}

func (hub *Hub) userJoinChannel(notification *Notification) {
	channelName := gjson.GetBytes(notification.Content, "data.channel").String()
	channel, ok := hub.channels[channelName]
	if !ok {
		channel = hub.newChannel(channelName)
	}
	channel.AddClient(notification.Client)
	notification.Client.AddChannel(channelName)
	channel.Broadcast(notification)

	hub.sendChannelInfo(channel, notification.Client)
}

func (hub *Hub) dealWithUserMessage(notification *Notification) {
	channelName := gjson.GetBytes(notification.Content, "data.channel").String()
	hub.broadcastOnChannel(channelName, notification)
}

func (hub *Hub) setClientUsername(notification *Notification) {
	username := gjson.GetBytes(notification.Content, "data.username").String()
	notification.Client.Username = username
}

func (hub *Hub) sendChannelInfo(channel *Channel, client *Client) {
	var users []string
	var notification ChannelInfoNotification
	for client := range channel.clients {
		users = append(users, client.Username)
	}

	notification.Type = "channelInfo"
	notification.Data.ChannelName = channel.Name
	notification.Data.Users = users

	bytes, err := json.Marshal(notification)
	if err != nil {
		panic(err)
	}

	//log.Println(string(bytes))

	client.send <- bytes

}