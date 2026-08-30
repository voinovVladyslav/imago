// Package logger logs every incoming message
package logger

import (
	"log"

	"imago/pkg/broker"
)

func Run() {
	config := broker.BrokerConfig{
		User:  "imago",
		Pass:  "strongpassword",
		Host:  "localhost",
		Port:  "5672",
		Vhost: "imago",
	}
	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatalf("can't connect to broker: %s", err)
	}
	msgs, err := b.Consume("logs", "#")
	if err != nil {
		log.Fatalf("can't consume from the queue %s", err)
	}

	go func() {
		for m := range msgs {
			log.Printf("key: %s message: %s\n", m.RoutingKey, m.Body)
			m.Ack(false)
		}
	}()

	var forever chan struct{}
	<-forever
}
