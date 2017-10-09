package main

import (
	"fmt"

	"github.com/valyala/gorpc"
)

func handler(clientAddr string, request interface{}) interface{} {
	in, ok := request.([]byte)
	if !ok {
		panic(fmt.Errorf("expected byte array got %T", request))
	}

	out := make([]byte, len(in))
	for i := range in {
		out[i] = 0x55
	}

	return out
}

func main() {
	fmt.Printf("staring server at %s\n", address)
	if err := gorpc.NewTCPServer(address, handler).Serve(); err != nil {
		panic(err)
	}
}
