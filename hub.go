package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type Hub struct {
	subsMap map[string][]*Client
	chanMap map[string]chan []byte
	lock    sync.Mutex
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

func (h *Hub) unsuscribe(client_id string, topic string) {
	h.lock.Lock()
	fmt.Println(h.subsMap)
	defer h.lock.Unlock()
	for index, client := range h.subsMap[topic] {
		if client.id == client_id {
			h.subsMap[topic] = append(h.subsMap[topic][:index], h.subsMap[topic][index+1:]...)
		}

	}
}

func (h *Hub) publish(w http.ResponseWriter, r *http.Request) {
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

	for _, client := range clients {
		if ok {
			client.conn.Write(client.ctx, websocket.MessageText, []byte(publish_msg.Message))
			// h.chanMap[publish_msg.Topic] <- []byte(publish_msg.Message)
			fmt.Println("published your message")
		}
	}
	w.Write([]byte("published your message to channel " + publish_msg.Topic + " with client id " + publish_msg.Cliend_id))

}

func (h *Hub) writeMessage(ctx context.Context, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)

}

func (h *Hub) listenIncomingMessage(ctx context.Context, c *websocket.Conn) error {

	for {
		_, msg, err := c.Reader(ctx)

		if err != nil {
			log.Fatal(err)
		}

		var wsMessage socketMsg
		// wsmsg_err := wsjson.Read(ctx, c, wsMessage)

		wsDecoder := json.NewDecoder(msg)
		decode_err := wsDecoder.Decode(&wsMessage)
		log.Println(wsMessage)
		if decode_err != nil {
			fmt.Println(decode_err)
		}

		switch wsMessage.Action {
		case "subscribe":
			newClient := Client{id: wsMessage.Cliend_id, conn: c, ctx: ctx}
			h.subsMap[wsMessage.Topic] = append(h.subsMap[wsMessage.Topic], &newClient)
			newClient.conn.Write(newClient.ctx, websocket.MessageText, []byte("added you to "+wsMessage.Topic))
			fmt.Println(h.subsMap[wsMessage.Topic])
		case "unsubscribe":
			for index, client := range h.subsMap[wsMessage.Topic] {
				if client.id == wsMessage.Cliend_id {
					h.subsMap[wsMessage.Topic] = append(h.subsMap[wsMessage.Topic][:index], h.subsMap[wsMessage.Topic][index+1:]...)
				}

			}
			fmt.Println(h.subsMap[wsMessage.Topic])
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
	}

	return h
}
