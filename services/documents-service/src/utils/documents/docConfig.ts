import * as compileDoc from '@services/index.js';
import * as schemas from '@events/index.js';

import paths from '@shared/constants/paths.js';

export type DocType = 'resume' | 'cover-letter';

export type ContentType = 'application/pdf' | 'text/plain';

interface DocConfig {
  pathFn: (uid: string, jobId: number) => string;
  jsonPathFn: (uid: string, docType: string) => string;
  generate: (
    docRequest: schemas.CompilationRequest
  ) => Promise<schemas.CompilationResult>;
  filename: (company: string, jobId: number) => string;
  contentType: ContentType;
}

export const docConfig: Record<DocType, DocConfig> = {
  resume: {
    pathFn: paths.paths.tempPdf,
    jsonPathFn: paths.paths.tempJson,
    generate: async (docRequest) => compileDoc.compileResume(docRequest),
    filename: (company, jobId) => `${company}_resume_${jobId}.pdf`,
    contentType: 'application/pdf' as const,
  },
  'cover-letter': {
    pathFn: paths.paths.tempPdf,
    jsonPathFn: paths.paths.tempJson,
    generate: async (docRequest) => compileDoc.compileCoverLetter(docRequest),
    filename: (company, uid) => `${company}_cover_letter_${uid}.pdf`,
    contentType: 'application/pdf',
  }
};
