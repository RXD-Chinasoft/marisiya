package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"html/template"
)

func main() {
	var upgrader = websocket.Upgrader{}
	var indexHandler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello" + time.Now().Format("2006-01-02 15:04:05"))
		
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}

	var homeTemplate = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprint(w, "home")
		homeTemplate.ExecuteTemplate(w, "home.html", 0)
	})
	http.ListenAndServe(":8000", nil)
}