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
	conn, err = amqp.Dial(rmqCredentials)
	if err != nil {
		return errors.New("Error de conexion: " + err.Error())
	}

	chann, err = conn.Channel()
	if err != nil {
		return errors.New("create channel error " + err.Error())
	}
	chann.Qos(1, 0, false)
	if err != nil {
		return errors.New("Error al abrir canal: " + err.Error())
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
		chann.Qos(1, 0, false)

	}

	if conn.IsClosed() {
		conn, _ = amqp.Dial(rmqCredentials)
		chann, _ = conn.Channel()
		chann.Qos(1, 0, false)
	}

}
