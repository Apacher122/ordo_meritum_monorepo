import { ResumeSectionNotFoundError } from "@shared/errors/resume_builder.errors.js";

export const sectionToLatexEnvMap: Record<
  string,
  "cvskills" | "cventries" | "cvletter" | "cvparagraph"
> = {
  summary: "cvparagraph",
  projects: "cventries",
  experiences: "cventries",
  skills: "cvskills",
  coverletter: "cvletter",
};

export const replaceSectionContent = (
  texContent: string,
  newContent: string[],
  sectionType: "cvskills" | "cventries" | "cvletter" | "cvparagraph"
) => {
  const environments = [
    { name: "cvskills", start: "\\begin{cvskills}", end: "\\end{cvskills}" },
    { name: "cventries", start: "\\begin{cventries}", end: "\\end{cventries}" },
    { name: "cvletter", start: "\\begin{cvletter}", end: "\\end{cvletter}" },
    { name: "cvparagraph", start: "\\begin{cvparagraph}", end: "\\end{cvparagraph}" }
  ];

  const env = environments.find((env) => env.name === sectionType);
  if (!env) {
    throw new ResumeSectionNotFoundError(
      `Unknown LateX environment for section: ${sectionType}`,
      { captureStackTrace: true }
    );
  }

  const { start, end } = env;
  const startIndex = texContent.indexOf(start);
  const endIndex = texContent.indexOf(end);

  if (startIndex === -1 || endIndex === -1) {
    throw new ResumeSectionNotFoundError(
      `No ${sectionType} environment found in the file.`,
      { captureStackTrace: true }
    );
  }

  const newEnvContent = `${start}\n${newContent.join("\n")}\n${end}`;
  return (
    texContent.slice(0, startIndex) +
    newEnvContent +
    texContent.slice(endIndex + end.length)
  );
};