package main

import (
	"net/http"
	"marisiya/db"
	. "marisiya/handlers"
)

func main() {
	db.AddFriend() //test
	http.HandleFunc("/ws", HandleWs)
	http.HandleFunc("/", HandleHomeByTemplate)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8000", nil)
}