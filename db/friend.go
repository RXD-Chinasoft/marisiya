package db

import (
	"strconv"
	"strings"
	// "context"
	"errors"  
	"database/sql"
	"log"
	. "marisiya/model"
	. "marisiya/protocal"
	"github.com/lib/pq"
	"fmt"
)



// api 1
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
		// some[cur].Friends = make([]int64, len(tmp))
		for _, t := range tmp {
			var has = false
			for _, m := range some[cur].Friends {
				if m == t {
					has = true
				}
			}
			if !has {
				some[cur].Friends = append(some[cur].Friends, t)
			}
		}
		// copy(some[cur].Friends, tmp)
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
	return true, nil

}

// api 2
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

	if len(hostInfo.Friends) < len(friends) - 1 {
		return 
	}
	for _, myFriend := range hostInfo.Friends {
		if isFriend = validRelation(myFriend, all, friends[1:]...);!isFriend {
			return
		}
	}
	return
}

// api 3
func FindCommonFriends(friends ...string) (common []string, err error) {
	if len(friends) <= 1 {
		err = errors.New("more than two friends required")
		return
	}
	some, err := GetSome(friends...)
	if err != nil {
		return
	}
	times := make(map[int64] int)
	common = []string{}
	log.Printf("common some %s \n", some)
	for _, v := range some {
		for _, m := range v.Friends {
			times[m] = times[m] + 1
			if times[m] == len(some) {
				var friend Friend
				friend, err = getFriendById(m)
				if err != nil {
					return
				}
				common = append(common, friend.Email)
			}
		}
	}

	log.Printf("common end %v \n", times)
	return
}

// api4
func Subscribe(requstor string, target string) (success bool, err error) {
	var hostInfo Friend
	hostInfo, err = GetInfoByEmail(target)
	if err != nil {
		return
	}
	var has = false
	for index, v := range hostInfo.SubscribMgr {
		if strings.Split(v, ",")[0] == requstor {
			hostInfo.SubscribMgr[index] = requstor +  "," + strconv.Itoa(1)
			has = true
			break
		}
	}
	if !has {
		hostInfo.SubscribMgr = append(hostInfo.SubscribMgr, requstor + "," + strconv.Itoa(1))
	}
	_, err = dbHandler.Exec("UPDATE friends set subscribMgr=$1 WHERE email=$2", pq.Array(hostInfo.SubscribMgr), target)
	if err != nil {
		err = errors.New("subscribe error")
		return
	}
	success = true
	return
}

// api5
func Block(requstor string, target string) (success bool, err error) {
	var hostInfo Friend
	hostInfo, err = GetInfoByEmail(target)
	if err != nil {
		return
	}
	var has = false
	for index, v := range hostInfo.SubscribMgr {
		if strings.Split(v, ",")[0] == requstor {
			hostInfo.SubscribMgr[index] = requstor +  "," + strconv.Itoa(0)
			has = true
			break
		}
	}
	if !has {
		hostInfo.SubscribMgr = append(hostInfo.SubscribMgr, requstor + "," + strconv.Itoa(0))
	}
	_, err = dbHandler.Exec("UPDATE friends set subscribMgr=$1 WHERE email=$2", pq.Array(hostInfo.SubscribMgr), target)
	if err != nil {
		err = errors.New("block error")
		return
	}
	success = true
	return
}

// helpers

func getFriendById(id int64) (friend Friend, err error) {
	friend = Friend{}
	row := dbHandler.QueryRow("SELECT * FROM friends WHERE id = $1", id)
	err = row.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends), pq.Array(&friend.SubscribMgr))
	switch {
	case err == sql.ErrNoRows:
		err = errors.New(fmt.Sprintf("not found %d", id))
	case err != nil:
		err = errors.New(fmt.Sprintf("interval error %s", err))
	}
	return
}

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
	err = row.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends), pq.Array(&friend.SubscribMgr))
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
		err := rows.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends), pq.Array(&friend.SubscribMgr))
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
		err := rows.Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends), pq.Array(&friend.SubscribMgr))
		if err != nil {
			log.Printf("scan friend error %s :", err)
			break
		}
		list = append(list, friend)

	}
	return list, err
}

func validRelation(hostFriend int64, allFriends []Friend, friends ...string) (pass bool) {
	for _, v := range allFriends {
		if hostFriend == v.Id { // find friend
			log.Printf("my friend detail %s, %v \n", v, friends)
			for _, email := range friends{
				if email == v.Email {
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
	err = dbHandler.QueryRow("SELECT * FROM friends WHERE id = $1", id).Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends), pq.Array(&friend.SubscribMgr))
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
	err = dbHandler.QueryRow("SELECT * FROM friends WHERE email = $1", email).Scan(&friend.Id, &friend.Email, pq.Array(&friend.Friends), pq.Array(&friend.SubscribMgr))
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