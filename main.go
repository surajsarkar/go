package main

import (
	"fmt"
	"net/http"
)

func main() {
	// log.SetFlags(0)

	h := newHub()

	go func() {

		publishServer := http.NewServeMux()
		publishServer.HandleFunc("/publish", h.publish)
		publishServerError := http.ListenAndServe(":9000", publishServer)
		if publishServerError != nil {
			fmt.Println("Error while spinning publish server")
		}
	}()
	go func() {

		server := http.NewServeMux()
		server.HandleFunc("/", h.actionHandler)

		socket_server_err := http.ListenAndServe(":8000", server)
		if socket_server_err != nil {
			fmt.Println("Error while spinning action server")
		}
	}()

}
