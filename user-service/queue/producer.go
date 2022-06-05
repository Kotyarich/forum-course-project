package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewProducer() *Producer {
	var p Producer

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	p.failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	p.failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"user-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	p.failOnError(err, "Failed to declare a queue")

	return &Producer{
		channel: ch,
		queue:   q,
	}
}

func (p *Producer) failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func (p *Producer) Produce(body string) {
	err := p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	p.failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
