/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	auth(c)

	for {
		data(c)

		time.Sleep(10*time.Second)

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
		DeviceId: "8",
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
