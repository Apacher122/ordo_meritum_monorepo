import * as kafka from "@kafka/index.js";
import * as services from "@services/index.js";

import { CompilationRequestSchema } from "@events/index.js";
import { logger } from "@shared/utils/logger.js";

export async function startDocumentWorker() {
  await kafka.consumer.connect();
  await kafka.producer.connect();
  await kafka.consumer.subscribe({ topic: kafka.Topics.LATEX_COMPILATION_REQUEST });

  await kafka.consumer.run({
    eachMessage: async ({ message }) => {
      if (!message.value) return;
      console.log("Received Kafka message:", message.value.toString());
      let request;
      try {
        request = CompilationRequestSchema.parse(JSON.parse(message.value.toString()));
      } catch (err) {
        logger.error("Invalid Kafka message:", err);
        return;
      }

      const resultPayload = await services.generateIfNeeded(request);
      logger.info("Generated document:", resultPayload);

      await kafka.producer.send({
        topic: kafka.Topics.LATEX_COMPILATION_RESULT,
        messages: [{ key: String(request.jobID), value: JSON.stringify(resultPayload) }],
      });
    },
  });
}
