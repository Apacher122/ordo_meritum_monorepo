/**
 * Formats a given section of data into a LaTeX section.
 * @param sectionType The type of section to format (e.g. "summary", "experiences", etc.).
 * @param sectionData The data to format for the LaTeX section.
 * @returns A string representing the formatted LaTeX section.
 * @throws Error If the section type is unknown.
 */
export const formatLatexSection =
  (sectionType: string) => (sectionData: any) => {
    switch (sectionType) {
      case "summary":
        return formatSummary(sectionData);
      case "experiences":
        return formatExperiences(sectionData);
      case "skills":
        return formatSkills(sectionData);
      case "projects":
        return formatProjects(sectionData);
      case "coverletter":
        return formatCoverLetter(sectionData);
      default:
        throw new Error(`Invalid section type: ${sectionType}`);
    }
  };

/**
 * Replaces special characters in a given string with their LaTeX equivalents.
 * @param text The string to format.
 * @returns A string with all special characters replaced with their LaTeX equivalents.
 */

export const formatTextForLatex = (text: string) => {
  if (!text) return "";

  const replacements: Record<string, string> = {
    "\\": "\\textbackslash{}",
    "%": "\\%",
    "#": "\\#",
    "&": "\\&",
    $: "\\$",
    "{": "\\{",
    "}": "\\}",
    "^": "\\^{}",
    _: "\\_",
    "~": "\\~{}",
    "\u00A0": " ",
    "\u2000": " ",
    "\u2001": " ",
    "\u2002": " ",
    "\u2003": " ",
    "\u2004": " ",
    "\u2005": " ",
    "\u2006": " ",
    "\u2007": " ",
    "\u2008": " ",
    "\u2009": " ",
    "\u200A": " ",
    "\u200B": " ",
    "\u200C": " ",
    "\u200D": " ",
    "\u202F": " ",
    "\u205F": " ",
    "\u3000": " ",
    "\u2018": "'",
    "\u2019": "'",
    "\u201A": "'",
    "\u201B": "'",
    "\u201C": '"',
    "\u201D": '"',
    "\u2013": "-",
    "\u2014": "-",
    "\u2015": "-",
  };

  const escapeForRegex = (char: string) =>
    char.replaceAll(/[-/\\^$*+?.()|[\]{}]/g, "\\$&");

  const regex = new RegExp(
    `[${Object.keys(replacements).map(escapeForRegex).join("")}]`,
    "g"
  );

  let result = text.replace(regex, (match) => replacements[match] || match);
  result = result.replaceAll(/\s+/g, " ").trim();

  return result;
};

/**
 * Format an array of summaries for LaTeX
 *
 * @param {any} data - An array of objects containing a sentence property
 * @returns {string} - A string of formatted summaries
 */
const formatSummary = (data: any) => {
  const summaryArray = Array.isArray(data) ? data : [data];

  return summaryArray
    .map(({ sentence }: { sentence: string }) => formatTextForLatex(sentence))
    .join(" ");
};

/**
 * Format an array of experiences for LaTeX
 *
 * @param {any} data - An object containing company, position, start, end, and bulletPoints properties
 * @returns {string} - A string of formatted experiences
 */
const formatExperiences = (data: any) => {
  const items = data.bulletPoints
    .map(
      ({ text }: { text: string }) => `    \\item {${formatTextForLatex(text)}}`
    )
    .join("\n");

  return `
\\cventry
  {${formatTextForLatex(data.company)}} % Organization
  {${formatTextForLatex(data.position)}} % Job title
  {} % Location
  {${formatTextForLatex(data.start)} - ${formatTextForLatex(
    data.end
  )}} % Date(s)
  {
    \\begin{cvitems}
${items}
    \\end{cvitems}
  }`;
};

/**
 * Formats an array of skills for LaTeX
 *
 * @param {any} data - An object containing a category and skill properties
 * @returns {string} - A string of formatted skills
 */
const formatSkills = (data: any) => {
  const skillList = data.skill
    .map((s: string) => formatTextForLatex(s))
    .join(", ");
  return `
\\cvskill
  {${formatTitle(data.category)}} % Category
  {${skillList}} % Skills`;
};

/**
 * Formats an array of projects for LaTeX
 *
 * @param {any} data - An object containing role, name, status, and bulletPoints properties
 * @returns {string} - A string of formatted projects
 */
const formatProjects = (data: any) => {
  const items = data.bulletPoints
    .map(
      ({ text }: { text: string }) => `    \\item {${formatTextForLatex(text)}}`
    )
    .join("\n");
  return `
\\cventry
  {${formatTextForLatex(data.role)}} % Role
  {${formatTextForLatex(data.name)}} % Event
  {} % Location
  {${formatTextForLatex(data.status)}} % Date(s)
  {
    \\begin{cvitems}
${items}
    \\end{cvitems}
  }`;
};

/**
 * Formats a cover letter for LaTeX
 *
 * @param {any} data - An object containing about, experience, and whatIBring properties
 * @returns {string} - A string of formatted cover letter sections
 */
const formatCoverLetter = (data: any) => `
\\lettersection{About}
${formatTextForLatex(data.about)}

\\lettersection{Experience}
${formatTextForLatex(data.experience)}

\\lettersection{What I Bring}
${formatTextForLatex(data.whatIBring)}
`;

/**
 * Format a title string for LaTeX
 * Replaces '%' with '\%', '#' with '\#', '&' with '\&', and '_' with ' '.
 * Converts the string to lowercase and capitalizes the first letter of each word.
 * @param {string} str - The title string to format
 * @returns {string} - The formatted title string
 */
const formatTitle = (str: string): string => {
  return str
    .replaceAll(
      /[%#&_]/g,
      (match) =>
        ({
          "%": "\\%",
          "#": "\\#",
          "&": "\\&",
          _: " ",
        }[match] || match)
    )
    .toLowerCase()
    .replaceAll(/\b\w/g, (s) => s.toUpperCase());
};
