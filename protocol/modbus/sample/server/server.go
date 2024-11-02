package main

import (
	"log"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
	serv := mbserver.NewServer()
	err := serv.ListenTCP("127.0.0.1:1502")
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer serv.Close()
	// Wait forever
	serv.RegisterFunctionHandler(1, func(s *mbserver.Server, f mbserver.Framer) ([]byte, *mbserver.Exception) {
		log.Printf("Function 1 called\n")
		return []byte{0x01, 0x02}, nil
	})

	for {
		time.Sleep(1 * time.Second)
	}
}