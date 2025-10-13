package kafka

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

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
