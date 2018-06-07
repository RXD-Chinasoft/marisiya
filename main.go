package main

import (
	"net/http"
	_ "marisiya/db"
	. "marisiya/handlers"
	. "marisiya/protocal"
)

func main() {
	messageChan := make(chan Message)
	http.HandleFunc("/ws", HandleWsByChan(messageChan))
	http.HandleFunc("/", HandleHomeByChan(messageChan))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/isFriend", HandleIsFriend)
	http.HandleFunc("/friends", HandleFriends)
	http.ListenAndServe(":8000", nil)
}