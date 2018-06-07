package handlers

import (
	"net/http"
	"log"
	"marisiya/db"
	"io/ioutil"
	"encoding/json"
)

func HandleIsFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s /n", res)
		} else {
			log.Printf("%s /n", res)
			db.IsFriend([]string{"a", "b"}...)
		}
		
	} else {
		// http.NotFoundHandler()
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
	}
}

func HandleFriends(w http.ResponseWriter, r *http.Request) {
	friends, err := db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		log.Printf("friends all r : %v \n", friends)
		if err1 := json.NewEncoder(w).Encode(friends); err1 != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		
	}
}