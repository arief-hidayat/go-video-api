package main

import "github.com/arief-hidayat/go-video-api/server"

func main() {
	s := server.NewServer()
	s.Start(":8000")
}
