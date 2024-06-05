package main

import (
	"nhclassen/golypus/server"
)

func main() {
	conf := server.ServerConfiguration{
		Address: "127.0.0.1",
		Port: "9999",
		Proto: "tcp",
	}

	s := server.NewServer(conf)

	s.Start()
}