package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitExchange     = "sat-stream"
	rabbitExchagneType = "topic"
)

func DoRabbitMQ() error {
	conn, err := amqp.Dial("amqp://user:password@host:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		rabbitExchange,     // name
		rabbitExchagneType, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	routing_key := "mapviewer.consume.success.#"

	err = ch.QueueBind(
		queue.Name,     // queue name
		routing_key,    // routing key
		rabbitExchange, // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")

	go func() {
		for d := range msgs {
			log.Printf(" [%s] %s", d.RoutingKey, d.Body)
		}
	}()
	var forever chan struct{}

	<-forever

	return nil
}
