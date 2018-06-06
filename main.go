package main

import (
	"net/http"
	"marisiya/db"
	. "marisiya/handlers"
)

func main() {
	db.AddFriend() //test
	http.HandleFunc("/ws", HandleWs)
	http.HandleFunc("/", HandleHome)
	http.ListenAndServe(":8000", nil)
}