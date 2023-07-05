package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/fetch"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"
)

const (
	rabbitExchange     = "sat-stream"
	rabbitExchagneType = "topic"
)

type RabbitMQNotification struct {
	Product   string `yaml:"product"`
	Timestamp string `yaml:"timestamp"`
	Project   string `yaml:"project"`
}

func DoRabbitMQ(outdir string, apihost string, sqldata types.SQLData) error {

	passw := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	url := fmt.Sprintf("amqp://user:%s@%s:5672/", passw, host)
	//conn, err := amqp.Dial("amqp://user:password@host:5672/")
	conn, err := amqp.Dial(url)
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
			//outdir := "./mapfiles/"
			notification := RabbitMQNotification{}

			//var data interface{}
			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = fetch.FetchProduct(notification.Project, notification.Product, outdir, apihost, sqldata)
			if err != nil {
				fmt.Println(err)
				continue
			}
			log.Printf(" [%s] %s", d.RoutingKey, d.Body)
		}
	}()
	var forever chan struct{}

	<-forever

	return nil
}
