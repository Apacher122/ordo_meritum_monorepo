import * as fs from "fs";
import * as path from "path";
import * as schemas from "@events/index.js";

import { formatTextForLatex } from "../latex/latex_formatters.js";

export const createHeader = async (
  uid: string,
  userInfo: schemas.UserInfo,
  education: schemas.EducationInfo,
  tempFolder: string,
  docType?: string,
  companyName?: string,
  jobTitle?: string
) => {
  const classFile = fs.readFileSync(path.join(tempFolder, "awesome-cv.cls"));
  const newClassFile = replaceVariables(classFile.toString(), { uid: uid });
  await fs.promises.writeFile(
    path.join(tempFolder, "awesome-cv.cls"),
    newClassFile
  );

  if (docType !== "resume") {
    const coverLetterInfoTemplate = path.join(
      tempFolder,
      "templates",
      "coverletter-template.tex"
    );

    const coverLetterInfo = fs.readFileSync(coverLetterInfoTemplate, "utf-8");
    const coverLetterWithPaths = replaceVariables(coverLetterInfo, {
      uid: uid,
      company: companyName,
      position: jobTitle,
    });

    const coverLetterInfoWithVariables = replaceVariables(
      coverLetterWithPaths,
      userInfo
    );


    await fs.promises.writeFile(
      path.join(tempFolder, "compiled", "coverletter.tex"),
      coverLetterInfoWithVariables
    );
    return;
  }

  const userInfoTemplate = path.join(
    tempFolder,
    "templates",
    "resume-template.tex"
  );
  const educationTemplate = path.join(
    tempFolder,
    "templates",
    "education-template.tex"
  );

  const resumeInfo = fs.readFileSync(userInfoTemplate, "utf-8");
  const resumeWithPaths = replaceVariables(resumeInfo, { uid: uid });
  const resumeInfoWithVariables = replaceVariables(resumeWithPaths, userInfo);

  const educationInfo = fs.readFileSync(educationTemplate, "utf-8");
  const educationInfoWithVariables = replaceVariables(educationInfo, education);

  await fs.promises.writeFile(
    path.join(tempFolder, "compiled", "resume.tex"),
    resumeInfoWithVariables
  );
  await fs.promises.writeFile(
    path.join(tempFolder, "compiled", "education.tex"),
    educationInfoWithVariables
  );
};



const replaceVariables = (
  template: string,
  data: Record<string, any>
): string => {
  let result = template;

  for (const key in data) {
    if (Object.prototype.hasOwnProperty.call(data, key)) {
      const regex = new RegExp(`<<${key}>>`, "g");
      result = result.replace(regex, formatTextForLatex(data[key] as string));
    }
  }
  return result;
};

// Load the basic information of the user into resume.tex
// export const loadUserInfoToLatex = async () => {
//   const resumeTemplate = await fs.promises.readFile(
//     paths.latex.resume.resumeTemplate,
//     "utf8"
//   );
//   const resumeInfo = Handlebars.compile(resumeTemplate)(infoStore.user_info);
//   await fs.promises.writeFile(paths.latex.resume.resume, resumeInfo);

//   const educationTemplate = await fs.promises.readFile(
//     paths.latex.resume.educationTemplate,
//     "utf8"
//   );

//   const educationInfo = Handlebars.compile(educationTemplate)(
//     infoStore.education_info
//   );
//   await fs.promises.writeFile(paths.latex.resume.education, educationInfo);
// };
