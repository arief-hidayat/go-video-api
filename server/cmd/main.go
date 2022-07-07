package main

import "go-video-api/server"

func main() {
	s := server.NewServer()
	s.Start(":8000")
}
