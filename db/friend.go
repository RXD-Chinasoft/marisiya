package db

import (
	"errors"  
	"database/sql"
	"log"
	. "marisiya/model"
	. "marisiya/protocal"
	"github.com/lib/pq"
	"fmt"
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
	err = row.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends))
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
	friend.Friends = []int64{1}
	_, err = dbHandler.Exec("INSERT INTO friends (id, email, friends) VALUES ($1, $2, $3)", friend.Id + 1, friend.Email, pq.Array(friend.Friends))
	if err != nil {
		log.Printf("insert error %s", err)
		return
	}
	return
}

func GetAll() ([]Friend, error) {
	list := []Friend{}
	rows, err := dbHandler.Query("SELECT * FROM friends;")
	if err != nil {
		log.Printf("get list error %s :", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		friend := Friend{}
		err := rows.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends))
		if err != nil {
			log.Printf("scan friend error %s :", err)
			break
		}
		list = append(list, friend)

	}
	return list, err
}

func IsFriend(friends ...string) (isFriend bool, err error) {
	if len(friends) < 1 {
		err = errors.New("more friends required")
		return
	}
	if len(friends) == 1 {
		isFriend = true
		return
	}
	log.Println(friends)
	var rfs []int64 //relation friends
	rfs, err = GetFriends(friends[0])
	if err != nil {
		return
	}
	if len(rfs) == 0 {
		log.Println("find no friends")
		return
	}
	log.Printf("relation friends %v :", rfs)
	var all []Friend // all friends
	all, err = GetAll()
	if err != nil {
		return
	}
	log.Printf("all friends %s :", all)

	for _, r := range rfs {
		bingo := false
		for _, v := range all {
			log.Printf("friend %s :", v)
			if r == v.Id {
				bingo = true
			}
		}
		if !bingo {
			isFriend = false;
			return
		}
	}
	isFriend = true;
	return
}

func GetFriends(mail string) (friends []int64, err error) {
	friend := Friend{}
	err = dbHandler.QueryRow("SELECT * FROM friends WHERE email = $1", mail).Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends))
	switch {
	case err == sql.ErrNoRows:
		log.Printf("%s do not exist", mail)
		err = errors.New(fmt.Sprintf("%s do not exist", mail))
		return
	case err != nil:
		err = errors.New(fmt.Sprintf("interval error :%s", err))
		log.Printf("interval error %s", err)
		return
	}
	friends = friend.Friends
	return
}