package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// log.SetFlags(0)

	var server_error_chan chan []byte

	h := newHub()

	go func() {

		publishServer := http.NewServeMux()
		publishServer.HandleFunc("/publish", h.publish)
		publishServerError := http.ListenAndServe(":9000", publishServer)
		if publishServerError != nil {
			server_error_chan <- []byte("error while starting publish server")
		}
	}()
	go func() {

		server := http.NewServeMux()
		server.HandleFunc("/", h.actionHandler)

		socket_server_err := http.ListenAndServe(":8000", server)
		if socket_server_err != nil {

			server_error_chan <- []byte("error while starting action server")
		}
	}()

	fmt.Println("started")
	msg := <-server_error_chan
	fmt.Println("see If staying")
	log.Println(msg)
}
