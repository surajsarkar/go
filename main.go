package main

import (
	"log"
	"net/http"
)

func main() {

	log.SetFlags(log.Default().Flags() | log.Llongfile)

	h := newHub()

	server := http.NewServeMux()
	server.HandleFunc("/publish", h.publish)
	server.HandleFunc("/ws", h.wsHandler)

	http.ListenAndServe(":8000", server)
}
