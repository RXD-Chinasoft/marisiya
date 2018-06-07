package db

import (
	"errors"  
	"database/sql"
	"log"
	. "marisiya/model"
	. "marisiya/protocal"
	"github.com/lib/pq"
)

func AddFriend(msg Message) (friend Friend, err error) {
	log.Println(msg)
	var v string
	var ok bool
	if v, ok = msg.Data.(string);!ok {
		err = errors.New("bad type")
		return
	}
	friend = Friend{}
	row := dbHandler.QueryRow("SELECT * FROM friends WHERE email = $1", v)
	err = row.Scan(&friend.Id, &friend.Email, &friend.Friends)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("not found %s", msg.Data)
	case friend.Email == v:
		log.Printf("user already exist %s", friend)
		err = errors.New("user already exist")
		return
	case err != nil:
		log.Printf("interval error %s", err)
		return
	}
	log.Println("begin insert")
	friend.Email = v
	_, err = dbHandler.Exec("INSERT INTO friends (id, email, friends) VALUES ($1, $2, $3)", friend.Id + 1, friend.Email, pq.Array(friend.Friends))
	if err != nil {
		log.Printf("insert error %s", err)
		return
	}
	return
}

func IsFriend(friends ...string) (hasFriend bool, err error) {
	if len(friends) <= 1 {
		return
	}
	log.Println(friends)
	rows, err := dbHandler.Query("SELECT * FROM friends;")
	if err != nil {
		log.Printf("find friends error %s :", err)
	}
	defer rows.Close()
	for rows.Next() {
		friend := &Friend{}
		err := rows.Scan(friend.Id, friend.Email, friend.Friends)
		if err != nil {
			// http.Error(w, http.StatusText(500), 500)
			break
		}

	}
	return
}