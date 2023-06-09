package main

import (
	"context"

	"nhooyr.io/websocket"
)

type Client struct {
	id   string
	conn *websocket.Conn
	ctx  context.Context
}
