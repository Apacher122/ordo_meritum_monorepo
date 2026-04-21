import { useEffect, useRef } from "react";

import { DocumentStatus } from "@/app/appProviders";
import { downloadDocument } from "../api/downloadDocument";

/**
 * Hook that downloads a document from the server and saves it to the local file
 * system.
 * @param {DocumentStatus | undefined} serverStatus - The status of the document
 *   on the server.
 * @param {number | null} jobId - The ID of the job.
 * @param {string} docType - The type of the document.
 * @param {string} companyName - The name of the company.
 * @param {string} jobTitle - The title of the job.
 * @param {() => Promise<void>} checkFile - A function that checks if a document
 *   with the given parameters exists.
 * @param {any} user - The user object.
 * @returns {void} Nothing is returned.
 */
export const useDocumentDownload = (
  serverStatus: DocumentStatus | undefined,
  jobId: number | null,
  docType: string,
  companyName: string,
  jobTitle: string,
  checkFile: () => Promise<void>,
  user: any,
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

    if (serverStatus.status !== "COMPLETED") {
      return;
    }

    downloadedJobsRef.current.add(jobKey);

    const downloadAndSave = async () => {
      try {
        console.log(`Downloading ${docType} for job ${jobId}...`);

        if (!user) throw new Error("User not available for authentication");

        const token = await user.getIdToken();
        const response = await downloadDocument(
          serverStatus.downloadUrl,
          serverStatus.changesUrl,
          token,
        );

        const pdfArrayBuffer = await response.pdf.arrayBuffer();
        const pdfPath = `${docType}/${companyName
          .toLowerCase()
          .replaceAll(" ", "_")}_${jobTitle
          .toLowerCase()
          .replaceAll(" ", "_")}_${jobId}_${docType}.pdf`;
        const jsonPath = `${docType}/${companyName
          .toLowerCase()
          .replaceAll(" ", "_")}_${jobTitle
          .toLowerCase()
          .replaceAll(" ", "_")}_${jobId}_${docType}.json`;

        await window.appAPI.files.saveFile(pdfPath, pdfArrayBuffer);
        await window.appAPI.files.saveJsonFile(jsonPath, response.jsonData);

        await checkFile();
        console.log(`Download complete: ${pdfPath}`);
      } catch (err) {
        console.error("Failed to download and save document:", err);
      }
    };

    downloadAndSave();
  }, [serverStatus, jobId, docType, companyName, jobTitle, checkFile, user]);
};
