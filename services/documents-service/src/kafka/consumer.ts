import { Kafka } from "kafkajs";

const kafkaBroker = process.env.KAFKA_BROKER_URL || "kafka:29092";

export const kafka = new Kafka({
  clientId: "document-service",
  brokers: [kafkaBroker],
});

export const consumer = kafka.consumer({ groupId: "latex-compilation-workers" });
