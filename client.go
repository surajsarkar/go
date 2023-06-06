package main

type Client struct {
	id   int
	msgs chan []byte
}
