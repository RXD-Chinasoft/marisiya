package handlers

import (
	"net/http"
	"log"
	"marisiya/db"
	"io/ioutil"
)

func HandleIsFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s /n", res)
		} else {
			log.Printf("%s /n", res)
			db.IsFriend()
		}
		
	} else {
		http.NotFoundHandler()
	}
}