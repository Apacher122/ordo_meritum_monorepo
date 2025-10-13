import * as fs from "fs";

import paths from "@shared/constants/paths.js";

export const initializeDocumentWorkspace = async (uid: string, docType: string) => {

  const tempFolder = paths.paths.tempDir(uid);
  const tempPdf = paths.paths.tempPdf(uid);
  const tempFolderCompiled = paths.latex.tempCompiled(uid);
  const tempJson = paths.paths.tempJson(uid, docType);

  fs.mkdirSync(tempFolder, { recursive: true });
  fs.mkdirSync(tempPdf, { recursive: true });
  fs.mkdirSync(tempFolderCompiled, { recursive: true });
  fs.mkdirSync(tempJson, { recursive: true });

  return { tempFolder, tempPdf, tempFolderCompiled, tempJson };
};
