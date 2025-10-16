import * as kafka from "@kafka/index.js";
import * as services from "@services/index.js";

import { CompilationRequestSchema } from "@events/index.js";
import { logger } from "@shared/utils/logger.js";

/**
 * Starts a document worker that consumes messages from the latex-compilation-requests Kafka topic,
 * generates a document based on the message payload, and produces a message to the latex-compilation-results topic.
 * The document worker will connect to the Kafka topic, subscribe to the latex-compilation-requests topic,
 * and run indefinitely until the process is exited.
 * Each message received from the topic will be parsed into a CompilationRequestSchema, and if the message is invalid,
 * an error will be logged and the message will be skipped.
 * If the message is valid, the document worker will call the generateIfNeeded function to generate a document based on the request.
 * The result of the generateIfNeeded function will be logged and a message will be sent to the latex-compilation-results topic.
 * The key of the message sent to the latex-compilation-results topic will be the jobID of the request, and the value will be the result payload as a JSON string.
 */
export async function startDocumentWorker() {
  await kafka.consumer.connect();
  await kafka.producer.connect();
  await kafka.consumer.subscribe({
    topic: kafka.Topics.LATEX_COMPILATION_REQUEST,
  });

  await kafka.consumer.run({
    eachMessage: async ({ message }) => {
      if (!message.value) return;
      console.log("Received Kafka message:", message.value.toString());
      let request;
      try {
        request = CompilationRequestSchema.parse(
          JSON.parse(message.value.toString())
        );
      } catch (err) {
        logger.error("Invalid Kafka message:", err);
        await kafka.producer.send({
          topic: kafka.Topics.LATEX_COMPILATION_RESULT,
          messages: [
            {
              key: String(request?.jobID ?? 0),
              value: JSON.stringify({error: "Invalid Kafka message"}),
            },
          ],
        });
        return
      }

      const resultPayload = await services.generateIfNeeded(request);
      logger.info("Generated document:", resultPayload);

      await kafka.producer.send({
        topic: kafka.Topics.LATEX_COMPILATION_RESULT,
        messages: [
          { key: String(request.jobID), value: JSON.stringify(resultPayload) },
        ],
      });
    },
  });
}
