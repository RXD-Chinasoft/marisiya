package handlers

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	. "marisiya/protocal"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
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

func HandleWsByChan(wsChan *WsChan) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade :: ", err)
			return
		}
		defer c.Close()
		wsChan.C = c
		for {
			message := &Message {}
			err := c.ReadJSON(message)
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			// err = c.WriteJSON(Message{Type:"welcome", Data:"welcome to websocket"})
			// switch message.Type {
			// case "new":
			// 	log.Printf("recv: %s", "new")
			// 	mchan <- *message
			// default:
			// 	err = c.WriteJSON(message)
			// }
			if mchan := wsChan.GroupChan[message.Type];mchan != nil {
				mchan <- *message
			} else {
				err = c.WriteJSON("bad request")
			}
				
			// err = c.WriteJSON(message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}