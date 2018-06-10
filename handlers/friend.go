package handlers

import (
	"net/http"
	"log"
	"marisiya/db"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type FriendsArray struct {
	Friends []int64 `json:"friends"`
}

type RequestForToBeFriends struct {
	Host string `json:"host"`
	Slaves []int64 `json:"friends"`
}

type Result struct {
	Success bool `json:"success"`
	Reseason string `json:"reseason"`
}

func HandleTobeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s \n", res)
		} else {
			param := RequestForToBeFriends{}
			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
			log.Printf("%+v \n", param)
			_, err = db.TobeFriend(param.Host, param.Slaves)
			result := Result{}
			if err != nil {
				result.Success = false
				result.Reseason = err.Error()
			} else {
				result.Success = true
				result.Reseason = ""
			}
			if err1 := json.NewEncoder(w).Encode(&result); err1 != nil {
				http.Error(w, http.StatusText(500), 500)
			}
		}

	}
}

func HandleIsFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s \n", res)
		} else {
			log.Printf("%s \n", string(res))
			param := FriendsArray{}
			// json.Unmarshal(res, &param)
			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
			log.Printf("%+v \n", param)
			var isFriend bool
			isFriend, err = db.IsFriend(param.Friends...)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
			} else {
				result := Result{}
				result.Success = isFriend
				if isFriend {
					result.Reseason = ""
				} else {
					result.Reseason = "they are not friends"
				}
				if err1 := json.NewEncoder(w).Encode(&result); err1 != nil {
					http.Error(w, http.StatusText(500), 500)
				}
			}
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