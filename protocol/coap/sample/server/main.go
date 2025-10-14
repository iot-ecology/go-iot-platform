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

import ( //导入相应的包

	"github.com/dustin/go-coap"
	"log"
	"net"
)

func handleA(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in handleA: path=%q: %#v from %v", m.Path(), m, a)

	log.Printf("Message = %s" , string(m.Payload))

	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("hello to you!"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)


		log.Printf("Transmitting from A %#v", res)
		return res
	}
	return nil
}

func handleB(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in handleB: path=%q: %#v from %v", m.Path(), m, a)
	if m.IsConfirmable() { //判断是否为CON数据
		res := &coap.Message{
			Type:      coap.Acknowledgement, //指定回复数据为ACK类型
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("good bye!"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain) //设置Option ContentFormat

		log.Printf("Transmitting from B %#v", res)
		return res
	}
	return nil
}


func main() {

	mux := coap.NewServeMux()
	mux.Handle("/a", coap.FuncHandler(handleA)) //创建 "/a"处理接口
	mux.Handle("/b", coap.FuncHandler(handleB)) //创建 "/b"处理接口

	log.Fatal(coap.ListenAndServe("udp", ":5683", mux)) //启动Server 端口为5683 这里为什么要用":5683"?搞不明白

}
