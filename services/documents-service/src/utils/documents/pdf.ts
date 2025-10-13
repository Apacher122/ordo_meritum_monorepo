import * as fs from "fs";

import { PDFDocument } from "pdf-lib";

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
