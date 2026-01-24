import { DocumentType, ResumeChanges } from "../types";
import { useCallback, useEffect, useMemo, useState } from "react";

import { downloadDocument } from "../api/downloadDocument";
import { generateDocument as generateDocumentApi } from "../api/generateDocument";
import { useAuth } from "../../auth/providers/AuthProvider";
import { useDocumentStatus } from "../providers/DocumentStatusProvider";
import { useSettings } from "../../settings/hooks/useSettings";
import { useUserInfo } from "@/features/user/hooks/useUserInfo";

const getSafeFilename = (
  jobId: number,
  docType: DocumentType,
  company: string,
  title: string
): string => {
  const sanitizePart = (input: string): string => {
    return (
      input
        .toLowerCase()
        .replaceAll(/[.,]/g, "")
        .replaceAll(/[/\s\\]+/g, '_')
        .replaceAll(/_+/g, "_")
        .replaceAll(/[\s/\\.]+/g, '_')
    );
  };

  const safeCompany = sanitizePart(company);
  const safeTitle = sanitizePart(title);

  const baseName = `${safeCompany}_${safeTitle}_${jobId}`;

  return `${baseName}_${docType}`;
};

/**
 * Hook that provides document management functionality.
 * @param {number | null} jobId - The ID of the job.
 * @param {string} companyName - The name of the company.
 * @param {string} jobTitle - The title of the job.
 * @param {DocumentType} docType - The type of the document.
 * @returns {object} An object containing the following properties:
 *   - displayStatus: The status of the document ("checking", "pending", "failed", "present", or "idle").
 *   - localPdfPath: The path to the locally saved PDF document if it exists.
 *   - localJsonData: The locally saved JSON data of the document if it exists.
 *   - generate: A function that generates the document using the API.
 *   - error: An error message if the document generation failed.
 *   - doesFileExist: A function that checks if a document with the given parameters exists.
 */
export const useDocumentManager = (
  jobId: number | null,
  companyName: string,
  jobTitle: string,
  docType: DocumentType
) => {
  const { user } = useAuth();
  const { settings } = useSettings();
  const { documentStatuses, addPendingDocument } = useDocumentStatus();
  const { userProfile } = useUserInfo();

  const [fileExists, setFileExists] = useState(false);
  const [localPdfPath, setLocalPdfPath] = useState<string | null>(null);
  const [localJsonData, setLocalJsonData] = useState<ResumeChanges | null>(
    null
  );
  const [isCheckingFile, setIsCheckingFile] = useState(true);
  const [apiError, setApiError] = useState<string | null>(null);

  const serverStatus = useMemo(
    () =>
      jobId ? documentStatuses.get(String(jobId))?.get(docType) : undefined,
    [documentStatuses, jobId, docType]
  );

  const checkFile = useCallback(async () => {
    if (!jobId) {
      setFileExists(false);
      setIsCheckingFile(false);
      return;
    }
    setIsCheckingFile(true);
    const pdfPath =
      getSafeFilename(jobId, docType, companyName, jobTitle) + ".pdf";
    const jsonPath =
      getSafeFilename(jobId, docType, companyName, jobTitle) + ".json";
    try {
      const exists = await window.appAPI.files.checkFileExists(pdfPath);
      setFileExists(exists);
      if (exists) {
        setLocalPdfPath(`../../public/pdfs/${pdfPath}`);
        const jsonResult = await window.appAPI.files.readJsonFile(jsonPath);
        setLocalJsonData(jsonResult.data);
      } else {
        setLocalPdfPath(null);
        setLocalJsonData(null);
      }
    } catch (err) {
      console.error("Error checking file:", err);
    } finally {
      setIsCheckingFile(false);
    }
  }, [jobId, docType, companyName, jobTitle]);

  useEffect(() => {
    checkFile();
  }, [checkFile]);

  useEffect(() => {
    if (!serverStatus || serverStatus.status !== "COMPLETED") return;
    if (!serverStatus.downloadUrl || !serverStatus.changesUrl) return;
    if (!jobId) return;
    const downloadAndSave = async () => {
      try {
        const token = await user?.getIdToken();
        const response = await downloadDocument(
          serverStatus.downloadUrl,
          serverStatus.changesUrl,
          token
        );
        const pdfBlob = response.pdf;
        const jsonData = response.jsonData;
        const arrayBuffer = await pdfBlob.arrayBuffer();
        const pdfPath =
          getSafeFilename(jobId, docType, companyName, jobTitle) + ".pdf";
        const jsonPath =
          getSafeFilename(jobId, docType, companyName, jobTitle) + ".json";
        await window.appAPI.files.saveFile(pdfPath, arrayBuffer);
        await window.appAPI.files.saveJsonFile(jsonPath, jsonData);
        await checkFile();
      } catch (err) {
        console.error("Failed to download and save file:", err);
        setApiError("Failed to download the new document.");
      }
    };
    downloadAndSave();
  }, [serverStatus, jobId, docType, companyName, jobTitle, checkFile, user]);

  const generate = useCallback(async () => {
    if (!jobId || !user || !settings || !userProfile) {
      setApiError("User, Job ID, Settings, or Profile are not loaded.");
      return;
    }

    setApiError(null);
    addPendingDocument(String(jobId), docType);

    try {
      const token = await user.getIdToken();
      const feature =
        docType === "resume" ? "resumeGeneration" : "coverLetterGeneration";
      const llmProvider = settings.featureAssignments[feature];

      await generateDocumentApi(
        docType,
        userProfile,
        jobId,
        llmProvider,
        settings,
        token
      );
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "An unknown error occurred.";
      console.error(
        `Failed to start ${docType} generation for job ${jobId}:`,
        errorMessage
      );
      setApiError(`API Error: ${errorMessage}`);
    }
  }, [jobId, docType, user, settings, addPendingDocument, userProfile]);

  const displayStatus = useMemo(() => {
    if (isCheckingFile) return "checking";
    if (serverStatus?.status === "PENDING") return "generating";
    if (serverStatus?.status === "FAILED" || apiError) return "failed";
    if (fileExists) return "present";
    return "idle";
  }, [isCheckingFile, serverStatus, fileExists]);

  const doesFileExist = useCallback(
    async (
      checkJobId: number,
      checkDocType: DocumentType,
      checkCompanyName: string,
      checkJobTitle: string
    ): Promise<boolean> => {
      const pdfPath =
        getSafeFilename(
          checkJobId,
          checkDocType,
          checkCompanyName,
          checkJobTitle
        ) + ".pdf";
      try {
        const exists = await window.appAPI.files.checkFileExists(pdfPath);
        return exists;
      } catch (err) {
        console.error("Error in doesFileExist check:", err);
        return false;
      }
    },
    []
  );

  return {
    displayStatus,
    localPdfPath,
    localJsonData,
    generate,
    error: apiError || serverStatus?.error,
    doesFileExist,
  };
};
