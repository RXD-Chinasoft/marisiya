package protocal

import (
	."github.com/gorilla/websocket"
)

const (
	KIND_HOME = "home"
	KIND_SUBSCRIBE = "subscribe"
)

type WsChan struct {
	GroupChan map[string]chan Message
	C *Conn
}

type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
	Cmd string `json:"cmd"`
}