package main

import (
	"encoding/json"
	"github.com/dustin/go-coap"
	"log"
	"time"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
}

func main() {
	c, err := coap.Dial("udp", "192.168.3.101:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	auth(c)

	for {
		data(c)

		time.Sleep(30*time.Second)

	}


}

func data(c*coap.Conn) {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   []byte("test"),
	}

	path := "/data"

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)



	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s", rv.Payload)
	}
}

func auth( c *coap.Conn) {
	auth := Auth{
		Username: "admin",
		Password: "admin",
		DeviceId: "1234567890",
	}
	marshal, _ := json.Marshal(auth)
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   marshal,
	}

	path := "/auth"

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s", rv.Payload)
	}
}
