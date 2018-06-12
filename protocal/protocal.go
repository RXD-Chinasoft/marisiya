package protocal

import (
	."github.com/gorilla/websocket"
)

const (
	KIND_HOME = "home"
	KIND_OTHERS = "others"
	CMD_NEW_FRIEND = "new_friend"
	CMD_PUBLISH = "publish"
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

type Publish struct {
	Emails []string `json:"emails"`
	Message string `json:"message"`
	Notes []string `json:"notes"`
}