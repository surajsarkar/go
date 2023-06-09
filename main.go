package main

import (
	"fmt"
	"net/http"
)

func main() {

	var server_spinning_err chan []byte

	h := newHub()

	go func() {

		publishServer := http.NewServeMux()
		publishServer.HandleFunc("/publish", h.publish)
		publishServerError := http.ListenAndServe(":9000", publishServer)
		if publishServerError != nil {
			server_spinning_err <- []byte("Error while spinning publish server")
		}
	}()
	go func() {

		server := http.NewServeMux()
		server.HandleFunc("/", h.actionHandler)

		socket_server_err := http.ListenAndServe(":8000", server)
		if socket_server_err != nil {
			server_spinning_err <- []byte("Error while spinning action server")
		}
	}()

	spinnng_err := <-server_spinning_err

	if spinnng_err != nil {
		fmt.Println(string(spinnng_err))
	}

}
