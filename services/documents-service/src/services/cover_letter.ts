import * as docs from "@utils/documents/index.js";
import * as fs from "fs";
import * as latex from "@utils/latex/index.js";
import * as schemas from "@events/index.js";

import dotenv from "dotenv";
import { exportLatex } from "./export.js";
import { logger } from "@shared/utils/logger.js";
import paths from "@shared/constants/paths.js";

dotenv.config();

export const compileCoverLetter = async (
  docRequest: schemas.CompilationRequest
): Promise<schemas.CompilationResult> => {
  const data = schemas.CoverLetterPayloadSchema.safeParse(
    docRequest.coverLetter
  );
  if (!data.success) {
    logger.error("Malformed cover letter request", data.error);
    return {
      user_id: docRequest.userID,
      job_id: docRequest.jobID,
      success: false,
      error: data.error.message,
    };
  }

  try {
    const { tempFolder, tempPdf, tempFolderCompiled, tempJson } =
      await docs.initializeDocumentWorkspace(
        docRequest.userID,
        docRequest.docType
      );
    const companyName = docs.companyNameToFile(docRequest.companyName);
    fs.cpSync(paths.latex.originalTemplate, tempFolder, { recursive: true });

    await docs.createHeader(
      docRequest.userID,
      docRequest.userInfo,
      docRequest.educationInfo,
      tempFolder,
      docRequest.docType,
      docRequest.coverLetter?.companyProperName,
      docRequest.coverLetter?.jobTitle
    );

    await latex.generateLatexSectionFile(
      "coverletter",
      docRequest.coverLetter?.body,
      tempFolder
    );

    const jsonFile = await docs.saveJson(
      docRequest.coverLetter,
      companyName,
      docRequest.jobID,
      tempJson,
      docRequest.docType
    );

    const pdfPath = await exportLatex({
      jobNameSuffix: "cover_letter",
      outputPath: tempPdf,
      compiledPdfPath: tempFolderCompiled,
      companyName: companyName,
      jobId: docRequest.jobID,
      docType: "coverletter",
    });

    return {
      user_id: docRequest.userID,
      job_id: docRequest.jobID,
      success: true,
      document_type: docRequest.docType,
      download_url: pdfPath,
      changes_url: jsonFile,
    };
  } catch (error) {
    logger.error("Failed to compile cover letter: " + (error as Error).message);
    return {
      user_id: docRequest.userID,
      job_id: docRequest.jobID,
      success: false,
      error: (error as Error).message,
    };
  }
};
