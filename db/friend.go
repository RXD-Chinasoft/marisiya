package db

import (
	"log"
	. "marisiya/model"
)

func AddFriend() {
	log.Println(dbHandler)
}

func IsFriend() (hasFriend bool, err error) {
	rows, err := dbHandler.Query("SELECT * FROM friends;")
	if err != nil {
		log.Printf("find friends error %s :", err)
	}
	defer rows.Close()
	for rows.Next() {
		friend := &Friend{}
		err := rows.Scan(friend.Id, friend.Email, friend.Friend)
		if err != nil {
			// http.Error(w, http.StatusText(500), 500)
			break
		}
		
	}
	return
}