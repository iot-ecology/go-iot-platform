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

// RMQ PACKAGE - "rmq"
import (
	"errors"
	"go.uber.org/zap"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rmqCredentials string = "amqp://guest:guest@localhost:5672"
	rmqContentType string = "application/json"
)

var conn *amqp.Connection
var chann *amqp.Channel

func ConnectToRMQ() (err error) {
	conn, err = amqp.Dial(genUrl(globalConfig.MQConfig))
	if err != nil {
		return errors.New("Error de conexion: " + err.Error())
	}

	chann, err = conn.Channel()
	if err != nil {
		return errors.New("create channel error " + err.Error())
	}
	err = chann.Qos(1, 0, false)

	if err != nil {
		return errors.New("qos setting error " + err.Error())
	}

	observeConnection()

	return nil
}

func observeConnection() {
	go func() {
		e := <-conn.NotifyClose(make(chan *amqp.Error))
		if e != nil {
			zap.S().Errorf("Conexion perdida: %+v\n", e)

		}
		zap.S().Errorf("Intentando reconectar con RMQ\n")

		closeActiveConnections()

		for err := ConnectToRMQ(); err != nil; err = ConnectToRMQ() {
			zap.S().Error(err)
			time.Sleep(5 * time.Second)
		}
	}()
}

// Can be also implemented in graceful shutdowns
func closeActiveConnections() {
	if chann.IsClosed() {
		channel, _ := conn.Channel()

		chann = channel
		err := chann.Qos(1, 0, false)
		if err != nil {
			zap.S().Errorf("qos setting error：%+v", err)

		}

	}

	if conn.IsClosed() {
		conn, _ = amqp.Dial(rmqCredentials)
		chann, _ = conn.Channel()
		err := chann.Qos(1, 0, false)
		if err != nil {
			zap.S().Errorf("qos setting error：%+v", err)

		}
	}

}
