package main

import "flag"

var (
	server  string
	total   int
	msgSize int
	limit   int
)

func init() {
	flag.StringVar(&server, "s", ":5555", "address:port of server")
	flag.IntVar(&total, "n", 5, "number of requests to send")
	flag.IntVar(&msgSize, "size", 60, "message size")
	flag.IntVar(&limit, "l", 100, "limit for messages to send ahead")

	flag.Parse()
}
