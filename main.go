package main

import (
	"net/http"
	_ "marisiya/db"
	. "marisiya/handlers"
	. "marisiya/protocal"
	// ."github.com/gorilla/websocket"
)

func main() {
	wsChan := WsChan{}
	wsChan.GroupChan = make(map[string]chan Message)
	wsChan.GroupChan[KIND_HOME] = make(chan Message)
	http.HandleFunc("/ws", HandleWsByChan(&wsChan))
	http.HandleFunc("/", HandleHomeByChan(&wsChan))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/isFriend", HandleIsFriend)
	http.HandleFunc("/friends", HandleFriends)
	// http.HandleFunc("/getFriends", HandleGetFriends)
	http.HandleFunc("/toBeFriends", HandleTobeFriends)
	http.HandleFunc("/retreiveCommonFriends", RetrieveCommonFriends)
	http.HandleFunc("/subscribe", HandleSubscribe)
	http.HandleFunc("/block", HandleBlock)
	http.ListenAndServe(":8000", nil)
}