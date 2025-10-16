export const formatLatexSection = (sectionType: string) => (sectionData: any) => {
  switch (sectionType) {
    case "summary": return formatSummary(sectionData);
    case "experiences": return formatExperiences(sectionData);
    case "skills": return formatSkills(sectionData);
    case "projects": return formatProjects(sectionData);
    case "coverletter": return formatCoverLetter(sectionData);
    default: throw new Error(`Invalid section type: ${sectionType}`);
  }
};


export const formatTextForLatex = (text: string) => {
  if (!text) return '';

  const replacements: Record<string, string> = {
    '\\': '\\textbackslash{}',
    '%': '\\%',
    '#': '\\#',
    '&': '\\&',
    '$': '\\$',
    '{': '\\{',
    '}': '\\}',
    '^': '\\^{}',
    '_': '\\_',
    '~': '\\~{}',
    '\u00A0': ' ',
    '\u2000': ' ',
    '\u2001': ' ',
    '\u2002': ' ',
    '\u2003': ' ',
    '\u2004': ' ',
    '\u2005': ' ',
    '\u2006': ' ',
    '\u2007': ' ',
    '\u2008': ' ',
    '\u2009': ' ',
    '\u200A': ' ',
    '\u200B': ' ',
    '\u200C': ' ',
    '\u200D': ' ',
    '\u202F': ' ',
    '\u205F': ' ',
    '\u3000': ' ',
    '\u2018': "'", 
    '\u2019': "'", 
    '\u201A': "'", 
    '\u201B': "'", 
    '\u201C': '"', 
    '\u201D': '"', 
    '\u2013': '-', 
    '\u2014': '-', 
    '\u2015': '-',
  };

  const escapeForRegex = (char: string) => char.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');

  const regex = new RegExp(
    `[${Object.keys(replacements).map(escapeForRegex).join('')}]`,
    'g'
  );

  let result = text.replace(regex, (match) => replacements[match] || match);
  result = result.replace(/\s+/g, ' ').trim();

  return result;
};

const formatSummary = (data: any) => {
    const summaryArray = Array.isArray(data) ? data : [data];

  return summaryArray
    .map(({ sentence }: { sentence: string }) => formatTextForLatex(sentence))
    .join(" ");

};

const formatExperiences = (data: any) => {
  const items = data.bulletPoints
    .map(({ text }: { text: string }) => `    \\item {${formatTextForLatex(text)}}`)
    .join("\n");

  return `
\\cventry
  {${formatTextForLatex(data.company)}} % Organization
  {${formatTextForLatex(data.position)}} % Job title
  {} % Location
  {${formatTextForLatex(data.start)} - ${formatTextForLatex(data.end)}} % Date(s)
  {
    \\begin{cvitems}
${items}
    \\end{cvitems}
  }`;
};

const formatSkills = (data: any) => {
  const skillList = data.skill.map((s: string) => formatTextForLatex(s)).join(", ");
  return `
\\cvskill
  {${formatTitle(data.category)}} % Category
  {${skillList}} % Skills`;
};

const formatProjects = (data: any) => {
  const items = data.bulletPoints.map(({ text }: { text: string }) => `    \\item {${formatTextForLatex(text)}}`).join("\n");
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

const formatCoverLetter = (data: any) => `
\\lettersection{About}
${formatTextForLatex(data.about)}

\\lettersection{Experience}
${formatTextForLatex(data.experience)}

\\lettersection{What I Bring}
${formatTextForLatex(data.whatIBring)}
`
;


const formatTitle = (str: string): string => {
    return str
    .replace(/[%#&_]/g, match => ({
        '%': '\\%',
        '#': '\\#',
        '&': '\\&',
        '_': ' '
      }[match] || match))
      .toLowerCase()
      .replace(/\b\w/g, s => s.toUpperCase());
};