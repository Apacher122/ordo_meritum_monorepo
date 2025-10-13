import * as path from 'path';

import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(__dirname, '../../..');

const resolveRoot = (...segments: string[]) => {
    return path.join(root, ...segments);
};
const latexRoot = resolveRoot('src', 'shared', 'templates', 'latex', 'original-template');
const SHARED_PDFS_ROOT = '/usr/src/app/shared_pdfs';
const SHARED_LOGS_ROOT = '/usr/src/app/logs/microservices/documents';

const outputPaths = {
    dir: resolveRoot('output'),
    tempPdf: (uid = '') =>path.join(SHARED_PDFS_ROOT, `${uid}`, 'pdf'),
    tempDir: (uid = '') =>path.join(SHARED_PDFS_ROOT, `${uid}`),
    tempJson: (uid = '', docType = '') =>path.join(SHARED_PDFS_ROOT, `${uid}`, 'json', `${docType}`),
}

const logPaths = {
    infoLogFile: path.join(SHARED_LOGS_ROOT, 'info.log'),
    errorLogFile: path.join(SHARED_LOGS_ROOT, 'error.log'),
}

const latexPaths = {
    class: path.join(latexRoot, 'awesome-cv.cls'),
    resume: {
        education: path.join(latexRoot, 'compiled', 'education.tex'),
        educationTemplate: path.join(latexRoot, 'templates', 'education-template.tex'),
        experiences: path.join(latexRoot, 'compiled', 'experience.tex'),
        experiencesTemplate: path.join(latexRoot, 'templates', 'experience-template.tex'),
        extracurriculars: path.join(latexRoot, 'compiled', 'extracurriculars.tex'),
        extracurricularsTemplate: path.join(latexRoot, 'templates', 'extracurriculars-template.tex'),
        honors: path.join(latexRoot, 'compiled', 'honors.tex'),
        honorsTemplate: path.join(latexRoot, 'templates', 'honors-template.tex'),
        projects: path.join(latexRoot, 'compiled', 'projects.tex'),
        projectsTemplate: path.join(latexRoot, 'templates', 'projects-template.tex'),
        skills: path.join(latexRoot, 'compiled', 'skills.tex'),
        skillsTemplate: path.join(latexRoot, 'templates', 'skills-template.tex'),
        summary: path.join(latexRoot, 'compiled', 'summary.tex'),
        summaryTemplate: path.join(latexRoot, 'compiled', 'summary-template.tex'),
        resume: path.join(latexRoot, 'resume.tex'),
        resumeTemplate: path.join(latexRoot, 'templates','resume-template.tex'),
    },
    coverLetter: {
        letter: path.join(latexRoot, 'coverletter.tex'),
        template: path.join(latexRoot, 'templates','coverletter-template.tex'),
    },
    cv: {
        cv: path.join(latexRoot, 'cv.tex'),
        cvTemplate: path.join(latexRoot, 'templates', 'cv-template.tex'),
    },
    originalTemplate: path.join(latexRoot),
    tempTemplate: (uid = '') =>path.join(SHARED_PDFS_ROOT, `${uid}`, 'templates'),
    tempCompiled: (uid = '') =>path.join(SHARED_PDFS_ROOT, `${uid}`, 'compiled'),
    template: (name: string) => path.join(latexRoot, 'templates', `${name}-template.tex`),
}

export default {
    root,
    resolveRoot,
    paths: {
        ...outputPaths,
        ...logPaths,
    },
    latex: latexPaths,
}

