package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func (h *Hub) actionHandler(w http.ResponseWriter, r *http.Request) {

	var susMsg suscriptionMsg

	parsed_body := json.NewDecoder(r.Body)
	decode_err := parsed_body.Decode(&susMsg)

	if decode_err != nil {
		w.Write([]byte("error while decoding body"))
	}

	if susMsg.Action == "suscribe" {
		h.suscribe(w, r, susMsg.Cliend_id, susMsg.Topic)
	} else if susMsg.Action == "unsuscribe" {
		h.unsuscribe(susMsg.Cliend_id)
	}

}

func (h *Hub) suscribe(w http.ResponseWriter, r *http.Request, client_id string, topic string) {

	c, err := websocket.Accept(w, r, nil)

	if err != nil {
		w.Write([]byte("error while upgrading connection"))
	}
	newClient := Client{id: client_id}
	_, ok := h.subsMap[topic]

	if ok {
		h.subsMap[topic] = append(h.subsMap[topic], &newClient)
	} else {

		h.subsMap[topic] = []*Client{&newClient}
	}
	listining_err := h.listenIncomingMessage()
	if listining_err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "connection closed")
}

func (h *Hub) unsuscribe(client_id string) {
	h.lock.Lock()
	fmt.Println(h.subsMap)
	defer h.lock.Unlock()
	for key, val := range h.subsMap {
		for client_index, client := range val {
			if client.id == client_id {
				h.subsMap[key] = append(h.subsMap[key][:client_index], h.subsMap[key][client_index+1:]...)
			}

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

	_, ok := h.subsMap[publish_msg.Topic]
	if ok {
		h.chanMap[publish_msg.Topic] <- []byte(publish_msg.Message)
	}

}

func (h *Hub) writeMessage(ctx context.Context, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)

}

func (h *Hub) listenIncomingMessage() error {

	for {
		select {
		case info_msg := <-h.chanMap["info"]:
			info_clients := h.subsMap["info"]
			for _, client := range info_clients {
				info_err := h.writeMessage(client.ctx, client.conn, info_msg)
				if info_err != nil {
					return info_err
				}
			}
		case broadcast_msg := <-h.chanMap["broadcast"]:
			for _, client := range h.subsMap["broadcast"] {
				brod_err := h.writeMessage(client.ctx, client.conn, broadcast_msg)
				if brod_err != nil {
					return brod_err
				}
			}

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
