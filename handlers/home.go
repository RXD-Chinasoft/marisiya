package handlers

import (
	"strconv"
	"strings"
	"fmt"
	// "encoding/json"
	"net/http"
	"log"
	"html/template"
	. "marisiya/protocal"
	"marisiya/db"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	homeTemplate := template.Must(template.New("").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
		<meta charset="utf-8">
		<script>  
		window.addEventListener("load", function(evt) {
			var output = document.getElementById("output");
			var input = document.getElementById("input");
			var ws;
			var print = function(message) {
				var d = document.createElement("div");
				d.innerHTML = message;
				output.appendChild(d);
			};
			document.getElementById("open").onclick = function(evt) {
				if (ws) {
					return false;
				}
				// ws = new WebSocket("{{.}}");
				ws = new WebSocket("ws://"+window.location.host+"/ws");
				ws.onopen = function(evt) {
					print("OPEN");
				}
				ws.onclose = function(evt) {
					print("CLOSE");
					ws = null;
				}
				ws.onmessage = function(evt) {
					print("RESPONSE: " + evt.data);
				}
				ws.onerror = function(evt) {
					print("ERROR: " + evt.data);
				}
				return false;
			};
			document.getElementById("send").onclick = function(evt) {
				if (!ws) {
					return false;
				}
				var message = JSON.stringify({type:"123", data: input.value})
                ws.send(message);
                print("SEND: " + message);
				return false;
			};
			document.getElementById("close").onclick = function(evt) {
				if (!ws) {
					return false;
				}
				ws.close();
				return false;
			};
		});
		</script>
		</head>
		<body>
		<table>
		<tr><td valign="top" width="50%">
		<p>Click "Open" to create a connection to the server, 
		"Send" to send a message to the server and "Close" to close the connection. 
		You can change the message and send multiple times.
		<p>
		<form>
		<button id="open">Open</button>
		<button id="close">Close</button>
		<p><input id="input" type="text" value="Hello world!">
		<button id="send">Send</button>
		</form>
		</td><td valign="top" width="50%">
		<div id="output"></div>
		</td></tr></table>
		</body>
		</html>
		`))
	err := homeTemplate.Execute(w, nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}

func HandleHomeByTemplate(w http.ResponseWriter, r *http.Request) {
	homeTemplate := template.Must(template.New("home.html").ParseFiles("templates/home.html"))
	err := homeTemplate.Execute(w, nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}

func HandleHomeByChan(wsChan *WsChan) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		homeTemplate := template.Must(template.New("home.html").ParseFiles("templates/home.html"))
		err := homeTemplate.Execute(w, nil)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
		go func() {
			for {
				message := <- wsChan.GroupChan[KIND_HOME]
				switch message.Cmd{
				case CMD_NEW_FRIEND:
					_, err = db.AddFriend(message) //test
					if err != nil {
						log.Println("Add friend fail: ", err)
					} else {
						wsChan.C.WriteJSON("add friend successfully")
					}
				case CMD_PUBLISH:
					if message.Data == nil || message.Data == "" {
						wsChan.C.WriteJSON("wrong data")
					} else {
						switch t := message.Data.(type) {
						default:
							fmt.Printf("unexpected type %T \n", t)
						}
						log.Printf("CMD_PUBLISH data: %v, %T\n", message.Data, message.Data)
						if publish, ok := message.Data.(map[string]interface{});!ok {
							wsChan.C.WriteJSON("wrong data")
						} else {
							friend, err := db.GetInfoByEmail(publish["email"].(string))
							if err != nil {
								wsChan.C.WriteJSON(err.Error())
							} else {
								log.Println("CMD_PUBLISH: ", friend)
								pbs := Publish{Emails:[]string{}}
								for _, v := range friend.SubscribMgr {
									if i, err := strconv.Atoi(strings.Split(v, ",")[1]);err != nil {
										wsChan.C.WriteJSON(err.Error())
									} else {
										if i == 1 {
											pbs.Emails = append(pbs.Emails, strings.Split(v, ",")[0])
										}
									}
								}
								pbs.Message = publish["message"].(string)
								wsChan.C.WriteJSON(pbs)
							}
						}
						
					}
					
				default:
					wsChan.C.WriteJSON("unkown cmd")
				}
				
			}
		}()
		
	}
}