package broker

import (
	"fmt"
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func Connect(user, pass, host, port string) (*ampq.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)

	conn, err := ampq.Dial(address)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(OrderCreatedPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, conn.Close
}
