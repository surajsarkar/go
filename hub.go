package main

type Hub struct {
	subscription_map map[string][]Client
	channel_map      map[string]chan string
}

func (h *Hub) Publish(topic string, message string, client Client) {
	user_present := false
	for _, c := range h.subscription_map[topic] {
		if c.id == client.id {
			user_present = true
		}
	}
	if !user_present {
		h.channel_map["issue"] <- "not suscribed"
	}

	h.channel_map[topic] <- message
}

func (h *Hub) Suscribe(topic string, client Client) {
	h.subscription_map[topic] = append(h.subscription_map[topic], client)
}

func (h *Hub) Unsuscribe(topic string, client Client) {
	for i, c := range h.subscription_map[topic] {
		if c.id == client.id {
			d := h.subscription_map[topic]
			h.subscription_map[topic] = append(h.subscription_map[topic][:i], d[i+1:]...)

		}
	}
}
