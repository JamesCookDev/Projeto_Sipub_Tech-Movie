package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
)

type MovieWriter interface {
	CreateMovie(ctx context.Context, title string, year int) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id string) error
}

type Consumer struct {
	service    MovieWriter
	url        string
	exchange   string
	exType     string
	queue      string
	rkCreated  string
	rkDeleted  string
}

func NewConsumer(s MovieWriter) *Consumer {
	return &Consumer{
		service:   s,
		url:       env("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		exchange:  env("RABBITMQ_EXCHANGE", "movies"),
		exType:    env("RABBITMQ_EXCHANGE_TYPE", "topic"),
		queue:     env("RABBITMQ_QUEUE", "movies.worker.q"),
		rkCreated: env("RABBITMQ_ROUTING_KEY_CREATED", "movie.created"),
		rkDeleted: env("RABBITMQ_ROUTING_KEY_DELETED", "movie.deleted"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	var conn *amqp.Connection
	var err error
	for i := 1; i <= 30; i++ {
		conn, err = amqp.Dial(c.url)
		if err == nil {
			break
		}
		log.Printf("[rabbitmq] tentativa %d/30: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	if err := ch.ExchangeDeclare(c.exchange, c.exType, true, false, false, false, nil); err != nil {
		_ = ch.Close(); _ = conn.Close(); return err
	}
	q, err := ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close(); _ = conn.Close(); return err
	}
	for _, rk := range []string{c.rkCreated, c.rkDeleted} {
		if err := ch.QueueBind(q.Name, rk, c.exchange, false, nil); err != nil {
			_ = ch.Close(); _ = conn.Close(); return err
		}
	}

	if err := ch.Qos(1, 0, false); err != nil {
		_ = ch.Close(); _ = conn.Close(); return err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		_ = ch.Close(); _ = conn.Close(); return err
	}

	go func() {
		for m := range msgs {
			if err := c.handle(m.RoutingKey, m.Body); err != nil {
				log.Printf("[consumer] erro processando (rk=%s): %v", m.RoutingKey, err)
				_ = m.Nack(false, true) // requeue
				continue
			}
			_ = m.Ack(false)
		}
	}()

	log.Printf("[consumer] ouvindo fila %s (rks: %s, %s)", c.queue, c.rkCreated, c.rkDeleted)

	<-ctx.Done()
	_ = ch.Close()
	_ = conn.Close()
	return nil
}

func (c *Consumer) handle(rk string, body []byte) error {
	var envelope struct {
		Action    string          `json:"action"`
		Data      json.RawMessage `json:"data"`
		Timestamp time.Time       `json:"timestamp"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}

	switch rk {
	case c.rkCreated:
		var req struct {
			Title string `json:"title"`
			Year  int32  `json:"year"`
		}
		if err := json.Unmarshal(envelope.Data, &req); err != nil {
			return err
		}
		_, err := c.service.CreateMovie(context.Background(), req.Title, int(req.Year))
		return err

	case c.rkDeleted:
		var d struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(envelope.Data, &d); err != nil {
			return err
		}
		return c.service.DeleteMovie(context.Background(), d.ID)
	}

	return nil
}

func env(k, fb string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return fb
}
