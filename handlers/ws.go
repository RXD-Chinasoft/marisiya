package handlers

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}
func init() {

}

func HandleWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade :: ", err)
		return
	}
	defer c.Close()
	for {
		message := &Message {}
		err := c.ReadJSON(message)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		// err = c.WriteJSON(Message{Type:"welcome", Data:"welcome to websocket"})
		err = c.WriteJSON(message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}