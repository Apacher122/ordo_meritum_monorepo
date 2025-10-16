import * as fs from "fs";

import { PDFDocument } from "pdf-lib";

/**
 * Takes a PDF file path and rewrites the PDF to a single page document.
 * This is useful for documents that are generated as multiple pages, but
 * need to be converted to a single page document for easier user
 * consumption.
 * @param {string} filePath - The path to the PDF file.
 * @returns {Promise<void>} - A promise that resolves when the PDF has been
 * rewritten to a single page document.
 */
export const forceSinglePagePDF = async (filePath: string): Promise<void> => {
  const existingPdfBytes = fs.readFileSync(filePath);
  const pdfDoc = await PDFDocument.load(new Uint8Array(existingPdfBytes));
  const pages = pdfDoc.getPages();
  
  for (let i = pages.length - 1; i >= 1; i--) {
    pdfDoc.removePage(i);
  }

  const pdfBytes = await pdfDoc.save();

  fs.writeFileSync(filePath, pdfBytes);
};
