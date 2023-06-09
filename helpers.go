package main

type suscriptionMsg struct {
	Action    string `json:"action"`
	Topic     string `json:"topic"`
	Cliend_id string `json:"client_id"`
}

type publishMsg struct {
	Cliend_id string `json:"client_id"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
}
