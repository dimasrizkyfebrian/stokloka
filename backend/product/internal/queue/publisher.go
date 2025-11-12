package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventPublisher adalah interface untuk publisher
type EventPublisher interface {
	Publish(routingKey string, data interface{}) error
}

// rabbitmqPublisher adalah implementasi nyata
type rabbitmqPublisher struct {
	ch *amqp.Channel
}

// NewRabbitMQPublisher membuat publisher baru
func NewRabbitMQPublisher(ch *amqp.Channel) EventPublisher {
	return &rabbitmqPublisher{ch: ch}
}

// Publish mengimplementasikan interface
func (p *rabbitmqPublisher) Publish(routingKey string, data interface{}) error {
	// Konversi data ke JSON
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshal data untuk event %s: %v", routingKey, err)
		return err
	}

	// Publish ke RabbitMQ
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Publishing event [Key: %s]", routingKey)

	err = p.ch.PublishWithContext(ctx,
		ExchangeName, // exchange (dari rabbitmq.go)
		routingKey,   // routing key (e.g., "product.created")
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("Gagal publish event: %v", err)
	}

	return nil
}
