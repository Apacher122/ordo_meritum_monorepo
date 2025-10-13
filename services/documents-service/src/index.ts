import { startDocumentWorker } from "@workers/documentWorker.js";

startDocumentWorker().catch((err: any) => {
  console.error("Document service crashed!", err);
  process.exit(1);
});
