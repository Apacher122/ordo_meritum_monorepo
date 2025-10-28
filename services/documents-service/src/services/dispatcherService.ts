import * as schemas from "@events/index.js";

import { docConfig } from "@utils/documents/index.js";

/**
 * Generates a document based on the compilation request if needed.
 *
 * This function is a dispatcher for different document generation functions.
 * It takes a compilation request as an argument and returns a compilation result.
 * The actual generation function is determined by the document type specified in the request.
 *
 * @param {schemas.CompilationRequest} docRequest Compilation request
 * @returns {Promise<schemas.CompilationResult>} Compilation result
 */
export const generateIfNeeded = async (
  docRequest: schemas.CompilationRequest
): Promise<schemas.CompilationResult> => {
  type DocType = keyof typeof docConfig;
  const { generate } = docConfig[docRequest.docType as DocType];

  const result = await generate(docRequest);

  return result;
};
