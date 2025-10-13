import * as path from "path";

import { forceSinglePagePDF } from "@shared/utils/documents/pdf/pdf.helpers.js";
import { logger } from "@shared/utils/logger.js";
import { spawn } from "child_process";
import { validatePath } from "@shared/utils/documents/file.helpers.js";

export const exportLatex = async ({
  jobNameSuffix,
  outputPath,
  compiledPdfPath,
  companyName,
  jobId,
  docType,
}: {
  jobNameSuffix: string;
  outputPath: string;
  compiledPdfPath: string;
  companyName: string;
  jobId: number;
  docType: string;
}): Promise<string> => {
  (compiledPdfPath);
  (outputPath);
  validatePath(outputPath);

  const latexFilePath = `${compiledPdfPath}/${docType}.tex`;
  const pdfPath = path.join(
    outputPath,
    `${companyName}_${jobNameSuffix}_${jobId}.pdf`
  );

  await executeLatex(
    companyName,
    jobNameSuffix,
    latexFilePath,
    jobId,
    outputPath
  );
  if (jobNameSuffix === "resume") {
    await forceSinglePagePDF(pdfPath);
  }

  return pdfPath;
};

const executeLatex = (
  companyName: string,
  jobNameSuffix: string,
  latexFilePath: string,
  id: number,
  outputPath: string
): Promise<void> => {
  return new Promise((resolve, reject) => {
    const latex = spawn("xelatex", [
      `--interaction=nonstopmode`,
      `-output-directory=${outputPath}`,
      `--jobname=${companyName}_${jobNameSuffix}_${id}`,
      latexFilePath,
    ]);

    latex.stdout.on("data", (data) => {
      logger.info("Sending LaTeX data: " + data.toString());
    });

    latex.stderr.on("data", (data) => {
      logger.error(`LaTeX Error: ${data.toString()}`);
    });

    latex.on("close", async (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`LaTeX process exited with code ${code}`));
      }
    });
  });
};
