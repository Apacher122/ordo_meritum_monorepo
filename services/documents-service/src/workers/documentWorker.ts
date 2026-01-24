import * as kafka from "@kafka/index.js";
import * as services from "@services/index.js";

import { CompilationRequestSchema } from "@events/index.js";
import { logger } from "@shared/utils/logger.js";

/**
 * Starts a document worker that consumes messages from the "latex-compilation-requests" Kafka topic.
 * Each message is parsed as a CompilationRequest and then passed to the generateIfNeeded function to generate a document.
 * The generated document is then sent as a message to the "latex-compilation-results" Kafka topic.
 * If an error occurs while parsing the message, an error message is sent to the "latex-compilation-results" topic instead.
 * The worker will continue to run until it is manually stopped.
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

      let request;
      try {
        request = CompilationRequestSchema.parse(
          JSON.parse(message.value.toString())
        );
        logger.info(`Received compilation request from user: ${request.userID}`);
      } catch (err) {
        logger.error("Invalid Kafka message:", err);
        console.log(JSON.stringify(request));
        await kafka.producer.send({
          topic: kafka.Topics.LATEX_COMPILATION_RESULT,
          messages: [
            {
              key: "error",
              value: JSON.stringify({ 
                error: "Invalid Kafka message",
                user_id: request?.userID ?? "",
                job_id: request?.jobID ?? "",
                success: false,
                document_type: request?.docType ?? "",
              }),
            },
          ],
        });
        return;
      }

      const resultPayload = await services.generateIfNeeded(request);
      logger.info(`Successfully generated document for user: ${request.userID}`);

      await kafka.producer.send({
        topic: kafka.Topics.LATEX_COMPILATION_RESULT,
        messages: [
          { key: String(request.jobID), value: JSON.stringify(resultPayload) },
        ],
      });
    },
  });
}
