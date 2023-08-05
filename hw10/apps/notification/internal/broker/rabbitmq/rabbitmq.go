package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"otus-microservices/notification/internal/broker"
	"otus-microservices/notification/internal/config"
	"time"
)
import amqp "github.com/rabbitmq/amqp091-go"

type Rabbitmq struct {
	Queue      string
	Host       string
	User       string
	Password   string
	Port       int
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func New(cfg config.RabbitmqConfig) *Rabbitmq {
	return &Rabbitmq{
		Queue:    cfg.Queue,
		Host:     cfg.Host,
		User:     cfg.User,
		Password: cfg.Password,
		Port:     cfg.Port,
	}
}

func (r *Rabbitmq) Connect() error {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
	connection, err := amqp.Dial(dsn)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		r.Queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	r.connection = connection
	r.channel = channel
	r.queue = &queue

	return nil
}

func (r *Rabbitmq) Send(message *broker.Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.channel.PublishWithContext(ctx,
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	return err
}

func (r *Rabbitmq) Receive() (*broker.Message, error) {
	messages, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return nil, err
	}

	delivery := <-messages
	body := delivery.Body

	result := broker.Message{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Rabbitmq) Close() error {
	err := r.channel.Close()
	if err != nil {
		return err
	}

	err = r.connection.Close()
	if err != nil {
		return err
	}

	return nil
}
