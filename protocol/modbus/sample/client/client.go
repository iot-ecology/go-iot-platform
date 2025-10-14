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
	"fmt"

	"github.com/goburrow/modbus"
)

func main() {
	handler := modbus.NewTCPClientHandler("localhost:1502")
	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()
	client := modbus.NewClient(handler)



	_, err = client.WriteMultipleRegisters(0, 5, []byte{0, 1, 0, 4, 0, 5, 1, 1, 1, 2})
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	results, err := client.ReadHoldingRegisters(0, 5)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("results %v\n", results)
}

