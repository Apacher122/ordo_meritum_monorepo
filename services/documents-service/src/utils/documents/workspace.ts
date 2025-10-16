import * as fs from "fs";

import paths from "@shared/constants/paths.js";

/**
 * Initializes a document workspace for a user.
 * This includes creating the temporary directories needed
 * for the document generation process.
 *
 * @param {string} uid - The user ID.
 * @param {string} docType - The type of document (e.g. resume, cover-letter).
 *
 * @returns {Promise<Object>} - A promise that resolves with an object containing the paths to the temporary directories.
 */
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
