import * as fs from "fs";
import * as path from "path";

import {
  formatLatexSection,
  replaceSectionContent,
  sectionToLatexEnvMap,
} from "./index.js";

export const generateLatexSectionFile = async (
  sectionName: string,
  data: any,
  tempFolder: string
) => {
  if (!data) return;

  const latexTemplatePath = path.join(
    tempFolder,
    "templates",
    `${sectionName}-template.tex`
  );
  const compiledLatexPath = path.join(
    tempFolder,
    "compiled",
    `${sectionName}.tex`
  );

  const originalLatexContent = await fs.promises.readFile(
    sectionName === "coverletter" ? compiledLatexPath : latexTemplatePath,
    "utf8"
  );

  const newContent = Array.isArray(data)
    ? data.map(formatLatexSection(sectionName))
    : [formatLatexSection(sectionName)(data)];

  const newLatexContent = replaceSectionContent(
    originalLatexContent,
    newContent,
    sectionToLatexEnvMap[sectionName]
  );

  await fs.promises.writeFile(compiledLatexPath, newLatexContent);
};
