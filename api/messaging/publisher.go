package messaging

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
	created  string
	deleted  string
}

func NewPublisher() (*Publisher, error) {
    url := env("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
    ex := env("RABBITMQ_EXCHANGE", "movies")
    exType := env("RABBITMQ_EXCHANGE_TYPE", "topic")
    rkCreated := env("RABBITMQ_ROUTING_KEY_CREATED", "movie.created")
    rkDeleted := env("RABBITMQ_ROUTING_KEY_DELETED", "movie.deleted")

    var conn *amqp.Connection
    var err error
    for i := 1; i <= 30; i++ {
        conn, err = amqp.Dial(url)
        if err == nil {
            break
        }
        log.Printf("[publisher] tentativa %d/30 conectar RabbitMQ: %v", i, err)
        time.Sleep(2 * time.Second)
    }
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        _ = conn.Close()
        return nil, err
    }
    if err := ch.ExchangeDeclare(ex, exType, true, false, false, false, nil); err != nil {
        _ = ch.Close()
        _ = conn.Close()
        return nil, err
    }
    _ = ch.Confirm(false)

    return &Publisher{
        conn: conn, ch: ch, exchange: ex,
        created: rkCreated, deleted: rkDeleted,
    }, nil
}


func (p *Publisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return p.ch.PublishWithContext(cctx,
		p.exchange, routingKey, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			DeliveryMode: amqp.Persistent,
		})
}

func (p *Publisher) RoutingKeyCreated() string { return p.created }
func (p *Publisher) RoutingKeyDeleted() string { return p.deleted }

func (p *Publisher) Close() {
	if p.ch != nil { _ = p.ch.Close() }
	if p.conn != nil { _ = p.conn.Close() }
}

func env(k, fb string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return fb
}

