import * as fs from "fs";
import * as path from "path";

/**
 * Saves a JSON file to disk with the given resume data.
 *
 * @param {any} resume - The resume data to be saved.
 * @param {string} company_name - The name of the company.
 * @param {number} jobId - The ID of the job.
 * @param {string} jsonPath - The path where the JSON file should be saved.
 * @param {string} docType - The type of the document (e.g. resume, cover-letter).
 * @returns {Promise<string>} - A promise that resolves with the path of the saved JSON file.
 */
export const saveJson = async (
  resume: any,
  company_name: string,
  jobId: number,
  jsonPath: string,
  docType: string
): Promise<string> => {
  const jsonFile = path.join(
    jsonPath,
    `${company_name}_${docType}_${jobId}.json`
  );

  fs.writeFileSync(jsonFile, JSON.stringify(resume, null, 2));

  return jsonFile;
};

/**
 * Converts a company name to a string that can be used as a file name.
 * Replaces spaces with underscores, removes non-alphanumeric characters, and
 * removes leading and trailing underscores.
 * @param {string} name - The company name to convert.
 * @returns {string} - The converted company name.
 */
export function companyNameToFile(name: string): string {
  if (!name) return "";
  let s = name.trim().toLowerCase();
  s = s.normalize("NFKD").replaceAll(/\p{M}/gu, "");
  s = s.replaceAll(/\s+/g, "_");
  s = s.replaceAll(/[^a-z0-9_]/g, "");
  s = s.replaceAll(/_+/g, "_").replaceAll(/(?:^_+)|(?:_+$)/g, "");
  return s;
}