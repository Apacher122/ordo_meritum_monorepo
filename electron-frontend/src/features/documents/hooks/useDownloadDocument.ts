import { useEffect, useRef } from 'react';

import { DocumentStatus } from '@/app/appProviders';
import { downloadDocument } from '../api/downloadDocument';

export const useDocumentDownload = (
  serverStatus: DocumentStatus | undefined,
  jobId: number | null,
  docType: string,
  companyName: string,
  jobTitle: string,
  checkFile: () => Promise<void>,
  user: any
) => {
  const downloadedJobsRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    if (
      !serverStatus ||
      !jobId ||
      !serverStatus.downloadUrl ||
      !serverStatus.changesUrl
    ) {
      return;
    }

    const jobKey = `${jobId}-${docType}`;
    if (downloadedJobsRef.current.has(jobKey)) {
      return;
    }

    if (serverStatus.status !== 'COMPLETED') {
      return;
    }

    downloadedJobsRef.current.add(jobKey);

    const downloadAndSave = async () => {
      try {
        console.log(`Downloading ${docType} for job ${jobId}...`);

        if (!user) throw new Error('User not available for authentication');

        const token = await user.getIdToken();
        const response = await downloadDocument(
          serverStatus.downloadUrl,
          serverStatus.changesUrl,
          token
        );

        const pdfArrayBuffer = await response.pdf.arrayBuffer();
        const pdfPath = `${docType}/${companyName.toLowerCase().replace(/ /g, '_')}_${jobTitle.toLowerCase().replace(/ /g, '_')}_${jobId}_${docType}.pdf`;
        const jsonPath = `${docType}/${companyName.toLowerCase().replace(/ /g, '_')}_${jobTitle.toLowerCase().replace(/ /g, '_')}_${jobId}_${docType}.json`;

        await window.appAPI.files.saveFile(pdfPath, pdfArrayBuffer);
        await window.appAPI.files.saveJsonFile(jsonPath, response.jsonData);

        await checkFile(); // update local state
        console.log(`Download complete: ${pdfPath}`);
      } catch (err) {
        console.error('Failed to download and save document:', err);
      }
    };

    downloadAndSave();
  }, [serverStatus, jobId, docType, companyName, jobTitle, checkFile, user]);
};
