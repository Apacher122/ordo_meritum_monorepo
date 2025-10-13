import * as schemas from "@events/index.js";

import { docConfig } from "@utils/documents/index.js";

export const generateIfNeeded = async (
  docRequest: schemas.CompilationRequest
): Promise<schemas.CompilationResult> => {
  type DocType = keyof typeof docConfig;
  const { generate } = docConfig[docRequest.docType as DocType];

  const result = await generate(docRequest);

  return result;
};
