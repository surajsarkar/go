package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

type Hub struct {
	subsMap map[string][]*Client
	chanMap map[string]chan []byte
	lock    sync.Mutex
	clients []Client
}

func (h *Hub) suscribe(w http.ResponseWriter, r *http.Request) {

	c, err := websocket.Accept(w, r, nil)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err != nil {
		w.Write([]byte("error while upgrading connection"))
	}
	c.Write(ctx, websocket.MessageText, []byte("got your message"))
	listining_err := h.listenIncomingMessage(ctx, c)
	if listining_err != nil {
		fmt.Println("error while listing")
		fmt.Println(listining_err)
	}
	defer c.Close(websocket.StatusInternalError, "connection closed")
}

func (h *Hub) publish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var publish_msg publishMsg
	publish_msg_decoder := json.NewDecoder(r.Body)

	err := publish_msg_decoder.Decode(&publish_msg)

	if err != nil {
		w.Write([]byte("error while reading body json"))
	}

	clients, ok := h.subsMap[publish_msg.Topic]

	if len(clients) == 0 {
		w.Write([]byte("no client present"))
	}

	if ok {
		for _, client := range clients {
			client.conn.Write(client.ctx, websocket.MessageText, []byte(publish_msg.Message))
			fmt.Println("published your message")
			w.Write([]byte("published your message to channel " + publish_msg.Topic + " with client id " + publish_msg.Cliend_id))
		}
	}

}

func (h *Hub) listenIncomingMessage(ctx context.Context, c *websocket.Conn) error {

	for {
		_, msg, err := c.Reader(ctx)

		if err != nil {
			log.Fatal(err)
		}

		var wsMessage socketMsg

		wsDecoder := json.NewDecoder(msg)
		err = wsDecoder.Decode(&wsMessage)
		log.Println(wsMessage)
		if err != nil {
			fmt.Println(err)
		}

		switch wsMessage.Action {
		case "register":
			newClient := Client{id: wsMessage.Cliend_id, conn: c, ctx: ctx}
			h.clients = append(h.clients, newClient)
			log.Println("registered " + newClient.id)
		case "subscribe":
			for _, client := range h.clients {
				if client.id == wsMessage.Cliend_id {
					h.subsMap[wsMessage.Topic] = append(h.subsMap[wsMessage.Topic], &client)
					client.conn.Write(client.ctx, websocket.MessageText, []byte("added you to "+wsMessage.Topic))
				}

			}
			fmt.Println(h.subsMap[wsMessage.Topic])
		case "unsubscribe":
			for index, client := range h.subsMap[wsMessage.Topic] {
				if client.id == wsMessage.Cliend_id {
					h.subsMap[wsMessage.Topic] = append(h.subsMap[wsMessage.Topic][:index], h.subsMap[wsMessage.Topic][index+1:]...)
				}

			}
			fmt.Println(h.subsMap[wsMessage.Topic])
		case "disconnect":
			for _, client := range h.clients {
				if client.id == wsMessage.Cliend_id {
					client.conn.Close(websocket.StatusGoingAway, "fullfilled my dreams")
					client.ctx.Done()
				}
			}
			log.Println("unregistered " + wsMessage.Cliend_id)
			// return nil
		case "publish":

			for _, client := range h.subsMap[wsMessage.Topic] {
				client.conn.Write(ctx, websocket.MessageText, []byte(wsMessage.Message))

			}
		default:
			continue

		}

	}

}

func newHub() *Hub {
	h := &Hub{
		subsMap: map[string][]*Client{"info": make([]*Client, 0)},
		chanMap: map[string]chan []byte{"info": make(chan []byte), "broadcast": make(chan []byte)},
		clients: make([]Client, 0),
	}

	return h
}
