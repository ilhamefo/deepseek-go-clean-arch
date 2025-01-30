package queue

import (
	"encoding/json"
	"event-registration/internal/core/domain"

	"github.com/streadway/amqp"
)

type EventQueue struct {
	conn *amqp.Connection
}

func NewEventQueue(conn *amqp.Connection) *EventQueue {
	return &EventQueue{conn: conn}
}

func (q *EventQueue) Publish(event *domain.Event) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",       // exchange
		"events", // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (q *EventQueue) Consume() (<-chan *domain.Event, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		"events", // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return nil, err
	}

	events := make(chan *domain.Event)
	go func() {
		for msg := range msgs {
			var event domain.Event
			json.Unmarshal(msg.Body, &event)
			events <- &event
		}
	}()

	return events, nil
}
