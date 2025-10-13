import { Kafka } from "kafkajs";

const kafkaBroker = process.env.KAFKA_BROKER_URL || "kafka:29092";

console.log(`Connecting to Kafka broker at: ${kafkaBroker}`);

export const kafka = new Kafka({
  clientId: "document-service",
  brokers: [kafkaBroker],
});

export const consumer = kafka.consumer({ groupId: "latex-compilation-workers" });

export const producer = kafka.producer();