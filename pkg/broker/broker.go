// Package broker used for inter-service message exchange
package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type BrokerConfig struct {
	User  string
	Pass  string
	Host  string
	Port  string
	Vhost string
}

const exchangeName string = "imago"

func NewBroker(c BrokerConfig) (*Broker, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", c.User, c.Pass, c.Host, c.Port, c.Vhost)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	b := &Broker{conn: conn, ch: ch}
	return b, nil
}

func (b *Broker) Consume(queueName, routingKey string) (<-chan amqp.Delivery, error) {
	q, err := b.ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = b.ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	msgs, err := b.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	return msgs, err
}

func (b *Broker) SendJSON(v any, key string) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return b.ch.PublishWithContext(
		ctx,
		exchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
