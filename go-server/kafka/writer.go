package kafka

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

// NewLatexWriter creates and configures a new kafka.Writer for producing messages
// to the "latex-compilation-requests" topic. It reads the Kafka broker URL from
// the "KAFKA_BROKER_URL" environment variable, falling back to a default if not set.
//
// This function is intended for use as an fx provider. It integrates with the fx.Lifecycle
// to manage the writer's connection, logging its initialization on start and closing
// it gracefully on stop.
func NewLatexWriter(lc fx.Lifecycle) *kafka.Writer {
	broker := os.Getenv("KAFKA_BROKER_URL")
	if broker == "" {
		broker = "kafka:29092"
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "latex-compilation-requests",
		Balancer: &kafka.LeastBytes{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().
				Str("service", "kafka-producer").
				Str("topic", writer.Topic).
				Msg("Kafka writer initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().
				Str("service", "kafka-producer").
				Msg("Closing Kafka writer...")
			return writer.Close()
		},
	})

	return writer
}
