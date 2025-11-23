package queue

import (
	"fmt"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Definisikan nama Exchange
const ExchangeName = "stokloka_topic_exchange"

// InitRabbitMQ membuat koneksi dan channel, lalu mendeklarasikan Exchange
func InitRabbitMQ(cfg *config.Config) (*amqp.Channel, func()) {
	// Buat DSN
	dsn := fmt.Sprintf("amqp://guest:guest@%s:%s/", cfg.RabbitHost, cfg.RabbitPort)

	// Buka Koneksi
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("Gagal konek ke RabbitMQ: %v", err)
	}

	// Buka Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Gagal membuka channel RabbitMQ: %v", err)
	}
	log.Println("Koneksi ke RabbitMQ dan Channel berhasil.")

	// Deklarasikan Exchange
	err = ch.ExchangeDeclare(
		ExchangeName, // nama
		"topic",      // tipe
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Gagal deklarasi Exchange: %v", err)
	}
	log.Printf("Exchange '%s' berhasil dideklarasikan.", ExchangeName)

	// Buat fungsi 'cleanup'
	cleanup := func() {
		ch.Close()
		conn.Close()
		log.Println("Koneksi RabbitMQ ditutup.")
	}

	return ch, cleanup
}
