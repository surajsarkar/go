package main

type Client struct {
	id   string
	msgs chan []byte
}
