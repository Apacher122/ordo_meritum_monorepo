package kafka

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/ordo_meritum/websocket"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

const serviceName = "kafka-consumer"

type consumer struct {
	reader *kafka.Reader
	hub    *websocket.Hub
}

type DocumentCompletionEvent struct {
	UserID       string `json:"user_id"`
	JobID        int    `json:"job_id"`
	Success      bool   `json:"success"`
	DocumentType string `json:"document_type"`
	DownloadURL  string `json:"download_url,omitempty"`
	ChangesURL   string `json:"changes_url,omitempty"`
	Error        string `json:"error,omitempty"`
}

func newConsumer(hub *websocket.Hub) *consumer {
	broker := os.Getenv("KAFKA_BROKER_URL")
	if broker == "" {
		broker = "kafka:29092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{broker},
		Topic:           "latex-compilation-results",
		GroupID:         "go-server-completion-consumers",
		MinBytes:        1,
		MaxBytes:        10e6,
		MaxWait:         10 * time.Second,
		ReadLagInterval: -1,
	})

	return &consumer{reader: reader, hub: hub}
}

func (c *consumer) start(ctx context.Context) {
	log.Info().Str("service", serviceName).Msg("Starting Kafka completion consumer...")
	defer func() {
		_ = c.reader.Close()
		log.Info().Str("service", serviceName).Msg("Kafka consumer stopped.")
	}()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Info().Msg("Kafka consumer context canceled, shutting down.")
				return
			}
			log.Error().Err(err).Msg("Kafka consumer error, retrying...")
			time.Sleep(2 * time.Second)
			continue
		}
		c.handleMessage(msg)
	}
}

func (c *consumer) handleMessage(msg kafka.Message) {
	var event DocumentCompletionEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal completion event")
		return
	}

	log.Info().
		Str("user_id", event.UserID).
		Str("job_id", strconv.Itoa(event.JobID)).
		Msg("Received completion event")

	c.broadcastEvent(&event, msg.Value)
}

func (c *consumer) broadcastEvent(event *DocumentCompletionEvent, rawMsg []byte) {
	if userClients, ok := c.hub.UserClients[event.UserID]; ok && len(userClients) > 0 {
		log.Info().
			Str("user_id", event.UserID).
			Int("client_count", len(userClients)).
			Msg("Broadcasting notification to connected clients")
		for client := range userClients {
			select {
			case client.Send <- rawMsg:
			default:
				close(client.Send)
				delete(userClients, client)
			}
		}
	} else {
		log.Warn().Str("user_id", event.UserID).Msg("No clients connected for user, cannot send notification")
	}
}

func RegisterCompletionConsumer(lc fx.Lifecycle, hub *websocket.Hub) {
	consumer := newConsumer(hub)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go consumer.start(context.Background())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Str("service", serviceName).Msg("Stopping Kafka completion consumer...")
			return consumer.reader.Close()
		},
	})
}
