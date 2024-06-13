package main

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

// Consumer holds all infromation
// about the RabbitMQ connection
// This setup does limit a consumer
// to one exchange. This should not be
// an issue. Having to connect to multiple
// exchanges means something else is
// structured improperly.
type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	consumerTag  string // Name that consumer identifies itself to the server with
	uri          string // uri of the rabbitmq server
	exchange     string // exchange that we will bind to
	exchangeType string // topic, direct, etc...
	bindingKey   string // routing key that we are using
}

// NewConsumer returns a Consumer struct
// that has been initialized properly
// essentially don't touch conn, channel, or
// done and you can create Consumer manually
func NewConsumer(
	consumerTag,
	uri,
	exchange,
	exchangeType,
	bindingKey string) *Consumer {
	return &Consumer{
		consumerTag:  consumerTag,
		uri:          uri,
		exchange:     exchange,
		exchangeType: exchangeType,
		bindingKey:   bindingKey,
		done:         make(chan error),
	}

}

// ReConnect is called in places where NotifyClose() channel is called
// wait 30 seconds before trying to reconnect. Any shorter amount of time
// will  likely destroy the error log while waiting for servers to come
// back online. This requires two parameters which is just to satisfy
// the AccounceQueue call and allows greater flexability
func (c *Consumer) ReConnect(queueName, bindingKey string) (<-chan amqp.Delivery, error) {
	time.Sleep(30 * time.Second)

	if err := c.Connect(); err != nil {
		zap.S().Errorf("Could not connect in reconnect call: %v", err.Error())
	}

	deliveries, err := c.AnnounceQueue(queueName, bindingKey)
	if err != nil {
		return deliveries, errors.New("Couldn't connect")
	}

	return deliveries, nil
}

// Connect to RabbitMQ server
func (c *Consumer) Connect() error {

	var err error

	zap.S().Infof("dialing %q", c.uri)
	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		zap.S().Infof("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Let Handle know it's not time to reconnect
		c.done <- errors.New("Channel Closed")
	}()

	zap.S().Infof("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	return nil
}

// AnnounceQueue sets the queue that will be listened to for this
// connection...
func (c *Consumer) AnnounceQueue(queueName, bindingKey string) (<-chan amqp.Delivery, error) {

	err := c.channel.Qos(100, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	zap.S().Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", c.consumerTag)
	deliveries, err := c.channel.Consume(
		queueName, // name
		"",        // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	return deliveries, nil
}

func (c *Consumer) Handle(
	d <-chan amqp.Delivery,
	fn func(<-chan amqp.Delivery),
	threads int,
	queue string,
	routingKey string) {

	var err error

	for {
		for i := 0; i < threads; i++ {
			go fn(d)
		}

		if <-c.done != nil {
			d, err = c.ReConnect(queue, routingKey)
			if err != nil {
				// Very likely chance of failing
				// should not cause worker to terminate
				zap.S().Fatalf("Reconnecting Error: %s", err)
			}
		}
		zap.S().Infof("Reconnected... possibly")
	}
}
