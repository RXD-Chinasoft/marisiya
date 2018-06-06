package protocal

type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
	Cmd string `json:"cmd"`
}