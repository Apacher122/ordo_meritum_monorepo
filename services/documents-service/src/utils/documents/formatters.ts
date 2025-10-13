import * as fs from "fs";
import * as path from "path";

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

export function companyNameToFile(name: string): string {
  if (!name) return "";
  let s = name.trim().toLowerCase();
  s = s.normalize("NFKD").replace(/\p{M}/gu, "");
  s = s.replace(/\s+/g, "_");
  s = s.replace(/[^a-z0-9_]/g, "");
  s = s.replace(/_+/g, "_").replace(/^_+|_+$/g, "");
  return s;
}