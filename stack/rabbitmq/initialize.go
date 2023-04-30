package rabbitmq_adapter

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(url string) *RabbitMQAdapter {
	conn, errConn := amqp.Dial(url)
	if errConn != nil {
		log.Fatalf("[Pubsage-RabbitMQ] Connection: %v", errConn)
	}

	ch, errCh := conn.Channel()
	if errCh != nil {
		log.Fatalf("[Pubsage-RabbitMQ] Channel: %v", errCh)
	}

	return &RabbitMQAdapter{conn: conn, ch: ch}
}

func NewConn(ch *amqp.Channel) *RabbitMQAdapter {
	return &RabbitMQAdapter{ch: ch}
}
