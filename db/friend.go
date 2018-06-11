package db

import (
	// "context"
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
	friend.Email = v
	friend.Friends = []int64{}
	log.Println("begin insert ", friend.Id)
	_, err = dbHandler.Exec("INSERT INTO friends (email, friends) VALUES ($1, $2)", friend.Email, pq.Array(friend.Friends))
	if err != nil {
		log.Printf("insert error %s", err)
		return
	}
	return
}

func TobeFriend(friends []string) (bool, error) {

	some, err := GetSome(friends...)
	if err != nil {
		return false, err
	}
	for _, f := range friends {
		tmp := []int64{}
		var cur = 0
		for i, v := range some {
			if v.Email != f {
				tmp = append(tmp, v.Id)
			} else {
				cur = i
			}
		}
		some[cur].Friends = make([]int64, len(tmp))
		copy(some[cur].Friends, tmp)
		log.Printf("some friends append %v , %v", tmp, some[cur])
	}
	// updates
	tx, err1 := dbHandler.Begin()
	if err1 != nil {
		return false, errors.New(fmt.Sprintf("begin update relation error %s", err1))
	}
	for _, v := range some {
		log.Printf("some friends update %s , %v", v.Email, v.Friends)
		dbHandler.Exec("UPDATE friends set friends=$1 WHERE email=$2", pq.Array(v.Friends), v.Email)
	}
	if err = tx.Commit();err != nil {
		return false, errors.New(fmt.Sprintf("commit update relation error %s", err))
	}

	// friend := Friend{}
	// row := dbHandler.QueryRow("SELECT * FROM friends WHERE email = $1", host)
	// err = row.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends))
	// switch {
	// case err == sql.ErrNoRows:
	// 	return false, errors.New(fmt.Sprintf("%s does not exist!", host))
	// case err != nil:
	// 	return false, errors.New(fmt.Sprintf("interval error :%s", err))
	// }
	// var all []Friend
	// all, err = GetAll()
	// if err != nil {
	// 	return false, err
	// }
	// relations := []int64{}
	// for _, v := range friends {
	// 	var exist = false
	// 	for _, m := range all {
	// 		if v == m.Email {
	// 			exist = true
	// 			relations = append(relations, m.Id)
	// 		}
	// 	}
	// 	if !exist {
	// 		return false, errors.New(fmt.Sprintf("%s does not exist!", v))
	// 	}
	// }
	// _, err = dbHandler.Exec("UPDATE friends set friends=$1 WHERE email=$2", pq.Array(relations), host)
	// if err != nil {
	// 	return false, errors.New(fmt.Sprintf("interval error :%s", err))
	// }
	return true, nil

}

func GetSome(emails ...string) ([]Friend, error) {
	if len(emails) == 0 {
		return nil, errors.New("emails number need more")
	}
	list := []Friend{}
	sentence := fmt.Sprintf("SELECT * FROM friends WHERE email = '%s'",emails[0])
	if ors := emails[1:];len(ors) > 0 {
		for _, email := range ors {
			sentence = sentence + " OR email = '" + email + "'"
		}
	}
	log.Println("sentence of to be friends ", sentence)
	rows, err := dbHandler.Query(sentence)
	if err != nil {
		log.Printf("get some list error %s :", err)
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
	// var rfs []int64 //relation friends
	// rfs, err = GetFriends(friends[0])
	var hostInfo Friend
	hostInfo, err = GetInfoByEmail(friends[0])
	if err != nil {
		return
	}
	if len(hostInfo.Friends) == 0 {
		log.Println("find no friends")
		return
	}
	log.Printf("relation friends %v \n", hostInfo.Friends)
	var all []Friend // all friends
	all, err = GetAll()
	if err != nil {
		return
	}
	log.Printf("all friends %s \n", all)

	// for _, r := range rfs {
	// 	bingo := false
	// 	for _, v := range all {
	// 		log.Printf("friend %s \n", v)
	// 		if r == v.Id {
	// 			bingo = true
	// 		}
	// 	}
	// 	if !bingo {
	// 		isFriend = false;
	// 		return
	// 	}
	// }
	for _, r := range hostInfo.Friends {
		if isFriend = validRelation(hostInfo.Id, r, all);!isFriend {
			return
		}
	}
	return
}

func validRelation(host int64, hostFriend int64, allFriends []Friend) (pass bool) {
	for _, v := range allFriends {
		if hostFriend == v.Id { // find friend
			log.Printf("my friend detail %s \n", v)
			for _, m := range v.Friends{
				if m == host {
					pass = true
					return
				}
			}
		}
	}
	return
}

func GetFriends(id int64) (friends []int64, err error) {
	friend := Friend{}
	err = dbHandler.QueryRow("SELECT * FROM friends WHERE id = $1", id).Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends))
	switch {
	case err == sql.ErrNoRows:
		log.Printf("%d do not exist", id)
		err = errors.New(fmt.Sprintf("%d do not exist", id))
		return
	case err != nil:
		err = errors.New(fmt.Sprintf("interval error :%s", err))
		log.Printf("interval error %d", err)
		return
	}
	friends = friend.Friends
	return
}

func GetInfoByEmail(email string) (friend Friend, err error) {
	friend = Friend{}
	err = dbHandler.QueryRow("SELECT * FROM friends WHERE email = $1", email).Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends))
	switch {
	case err == sql.ErrNoRows:
		log.Printf("%s do not exist", email)
		err = errors.New(fmt.Sprintf("%s do not exist", email))
		return
	case err != nil:
		err = errors.New(fmt.Sprintf("interval error :%s", err))
		log.Printf("interval error %d", err)
		return
	}
	// friends = friend.Friends
	return
}

func GetFriendsName(mail string) (friends []string, err error) {
	var hostInfo Friend
	// var rfs []int64 //relation friends
	hostInfo, err = GetInfoByEmail(mail)
	if err != nil {
		return
	}
	if len(hostInfo.Friends) == 0 {
		log.Println("find no friends")
		return
	}
	log.Printf("relation friends %v :", hostInfo.Friends)
	var all []Friend // all friends
	all, err = GetAll()
	if err != nil {
		return
	}
	log.Printf("all friends %s :", all)

	friends = []string{}
	for _, r := range hostInfo.Friends {
		for _, v := range all {
			log.Printf("friend %s :", v)
			if r == v.Id {
				friends = append(friends, v.Email)
			}
		}
	}
	return
}