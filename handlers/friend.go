package handlers

import (
	"net/http"
	"log"
	"marisiya/db"
	"io/ioutil"
	"encoding/json"
	"strings"
	// . "marisiya/protocal"
)

type FriendsArray struct {
	Friends []string `json:"friends"`
}

type RequestForToBeFriends struct {
	// Host string `json:"host"`
	Slaves []string `json:"friends"`
}

type RequestForSubscribe struct {
	Requestor string `json:"requestor"`
	Target string `json:"target"`
}

type RequestForSubscribeMention struct {
	Sender string `json:"sender"`
	Text string `json:"text"`
}

type Result struct {
	Success bool `json:"success"`
	Reseason string `json:"reseason"`
}

type PostEmail struct {
	Email string `json:"email"`
}

type CommonFriendsResponse struct {
	Success bool `json:"success"`
	Friends []string `json:"friends"`
	Count int `json:"count"`
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
			_, err = db.TobeFriend(param.Slaves)
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

// func HandleGetFriends(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		res, err := ioutil.ReadAll(r.Body)
// 		r.Body.Close()
// 		if err != nil {
// 			log.Printf("err %s \n", res)
// 			http.Error(w, http.StatusText(http.StatusBadRequest), 400)
// 		} else {
// 			log.Printf("%s \n", string(res))
// 			param := PostEmail{}
// 			// json.Unmarshal(res, &param)
// 			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
// 			log.Printf("%+v \n", param)

// 			friends, err := db.GetFriendsName(param.Email)
// 			if err != nil {
// 				http.Error(w, err.Error(), 500)
// 			} else {
// 				result := FindFriendsResult{}
// 				result.Success = true
// 				result.Friends = friends
// 				if err1 := json.NewEncoder(w).Encode(&result); err1 != nil {
// 					http.Error(w, http.StatusText(500), 500)
// 				}
// 				log.Printf("friends all name : %v \n", friends)
// 			}
// 		}
		
// 	} else {
// 		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
// 	}
// }

func RetrieveCommonFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s \n", res)
			http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		} else {
			log.Printf("%s \n", string(res))
			param := FriendsArray{}
			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
			log.Printf("%+v \n", param)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), 400)
			}
			friends, err := db.FindCommonFriends(param.Friends ...)
			if err != nil {
				http.Error(w, err.Error(), 500)
			} else {
				result := CommonFriendsResponse{}
				result.Success = true
				result.Friends = friends
				result.Count = len(friends)
				if err1 := json.NewEncoder(w).Encode(&result); err1 != nil {
					http.Error(w, http.StatusText(500), 500)
				}
				log.Printf("friends all name : %v \n", friends)
			}
		}
		
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
	}
}

func HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s \n", res)
			http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		} else {
			log.Printf("%s \n", string(res))
			param := RequestForSubscribe{}
			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
			log.Printf("%+v \n", param)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), 400)
			}
			success, err1 := db.Subscribe(param.Requestor, param.Target)
			result := Result{}
			if err1 != nil {
				result.Reseason = err1.Error()
				result.Success = false
			} else {
				result.Reseason = ""
				result.Success = success
			}
			if err = json.NewEncoder(w).Encode(&result);err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			}
			
		}
		
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
	}
}

func HandleBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s \n", res)
			http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		} else {
			log.Printf("%s \n", string(res))
			param := RequestForSubscribe{}
			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
			log.Printf("%+v \n", param)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), 400)
			}
			success, err1 := db.Block(param.Requestor, param.Target)
			result := Result{}
			if err1 != nil {
				result.Reseason = err1.Error()
				result.Success = false
			} else {
				result.Reseason = ""
				result.Success = success
			}
			if err = json.NewEncoder(w).Encode(&result);err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			}
		}
		
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
	}
}

func HandleSubMention(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Printf("err %s \n", res)
			http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		} else {
			log.Printf("%s \n", string(res))
			param := RequestForSubscribeMention{}
			err = json.NewDecoder(strings.NewReader(string(res))).Decode(&param)
			log.Printf("%+v \n", param)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), 400)
			}
			strs := strings.Split(param.Text, " ")
			if len(strs) != 2 {
				http.Error(w, http.StatusText(http.StatusBadRequest), 400)
				return
			}
			results, err1 := db.FindMentionedSubedEmails(param.Sender, strs[1])
			type MentionResult struct {
				Success bool `json:"success"`
				Recipients []string `json:"recipients"`
			}
			result := MentionResult{}
			if err1 != nil {
				result.Recipients = []string{}
				result.Success = false
			} else {
				result.Recipients = results
				result.Success = true
			}
			if err = json.NewEncoder(w).Encode(&result);err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			}
		}
		
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
	}
}