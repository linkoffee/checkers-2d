package main

import (
	"time"
)

type MessageType int

const (
	Success MessageType = iota
	Error
	Info
)

type Message struct {
	text         string
	messageType  MessageType
	showUntil    int64
}

func (g *Game) addMessage(text string, messageType MessageType, duration time.Duration) {
	expirationTime := time.Now().Add(duration).UnixNano()
	g.messages = append(g.messages, Message{text: text, messageType: messageType, showUntil: expirationTime})
}

