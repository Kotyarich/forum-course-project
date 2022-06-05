package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp.Channel
	tQueue  amqp.Queue
	vQueue  amqp.Queue
}

func NewProducer() *Producer {
	var p Producer

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	p.failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	p.failOnError(err, "Failed to open a channel")

	threadQ, err := ch.QueueDeclare(
		"thread-created", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	p.failOnError(err, "Failed to declare a queue")

	voteQ, err := ch.QueueDeclare(
		"vote-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	p.failOnError(err, "Failed to declare a queue")

	return &Producer{
		channel: ch,
		tQueue:  threadQ,
		vQueue:  voteQ,
	}
}

func (p *Producer) failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func (p *Producer) ProduceNewThread(body string) {
	err := p.channel.Publish(
		"",            // exchange
		p.tQueue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	p.failOnError(err, "Failed to publish thread message")
	log.Printf(" [x] thread sent %s\n", body)
}

func (p *Producer) ProduceNewVote(body string) {
	err := p.channel.Publish(
		"",            // exchange
		p.vQueue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	p.failOnError(err, "Failed to publish vote message")
	log.Printf(" [x] vote sent %s\n", body)
}
