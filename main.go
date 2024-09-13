package main

import (
	"nhclassen/golypus/server"
)

var (
	address string
	port    string
	proto   string
	static  string
)

func init() {
	address = "127.0.0.1"
	port = "80"
	proto = "tcp"
	static = "tcp"
}

func main() {
	conf := server.ServerConfiguration{
		Address: address,
		Port:    port,
		Proto:   proto,
		Static:  static,
	}

	s := server.NewServer(conf)

	s.Start()
}
