package data

import (
	"context"
	"kratos_single/internal/conf"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewMQ(c *conf.Data) (*MQ, error) {

	conn, err := amqp.Dial(c.Rabbitmq.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &MQ{
		conn: conn,
		ch:   ch,
	}, nil
}

// Publish 发消息
func (m *MQ) Publish(ctx context.Context, queue string, body string) error {
	_, err := m.ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,
	)
	if err != nil {
		return err
	}

	return m.ch.PublishWithContext(
		ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
}


// Consume 消费者
func (m *MQ) Consume(queue string) error {
	_, err := m.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := m.ch.Consume(
		queue,
		"",
		true, // auto ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Println("收到消息:", string(msg.Body))
		}
	}()

	return nil
}

// Close
func (m *MQ) Close() {
	if m.ch != nil {
		m.ch.Close()
	}
	if m.conn != nil {
		m.conn.Close()
	}
}