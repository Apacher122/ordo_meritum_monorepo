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

type DocumentCompletionEvent struct {
	UserID       string `json:"user_id"`
	JobID        int    `json:"job_id"`
	Success      bool   `json:"success"`
	DocumentType string `json:"document_type"`
	DownloadURL  string `json:"download_url,omitempty"`
	ChangesURL   string `json:"changes_url,omitempty"`
	Error        string `json:"error,omitempty"`
}

func RegisterCompletionConsumer(lc fx.Lifecycle, hub *websocket.Hub) {
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

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Str("service", serviceName).Msg("Starting Kafka completion consumer...")
			go func() {
				defer func() {
					_ = reader.Close()
					log.Info().Str("service", serviceName).Msg("Kafka consumer stopped.")
				}()
				ctx := context.Background()
				for {
					msg, err := reader.ReadMessage(ctx)
					if err != nil {
						if err == context.DeadlineExceeded {
							continue
						}
						if err == context.Canceled {
							log.Info().Msg("Kafka consumer context canceled, shutting down.")
							return
						}
						log.Error().Err(err).Msg("Kafka consumer error, retrying...")
						time.Sleep(2 * time.Second)
						continue
					}

					var event DocumentCompletionEvent
					if err := json.Unmarshal(msg.Value, &event); err != nil {
						log.Error().Err(err).Msg("Failed to unmarshal completion event")
						continue
					}

					log.Info().
						Str("user_id", event.UserID).
						Str("job_id", strconv.Itoa(event.JobID)).
						Msg("Received completion event")

					if userClients, ok := hub.UserClients[event.UserID]; ok {
						log.Info().
							Str("user_id", event.UserID).
							Int("client_count", len(userClients)).
							Msg("Broadcasting notification to connected clients")
						for client := range userClients {
							select {
							case client.Send <- msg.Value:
							default:
								close(client.Send)
								delete(userClients, client)
							}
						}
					} else {
						log.Warn().Str("user_id", event.UserID).Msg("No clients connected for user, cannot send notification")
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Str("service", serviceName).Msg("Stopping Kafka completion consumer...")
			return reader.Close()
		},
	})
}
