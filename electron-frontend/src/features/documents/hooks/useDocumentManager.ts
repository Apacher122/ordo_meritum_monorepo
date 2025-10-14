import { DocumentType, ResumeChanges } from "../types";
import { useCallback, useEffect, useMemo, useState } from "react";

import { downloadDocument } from "../api/downloadDocument";
import { generateDocument as generateDocumentApi } from "../api/generateDocument";
import { useAuth } from "../../auth/providers/AuthProvider";
import { useDocumentStatus } from "../providers/DocumentStatusProvider";
import { useSettings } from "../../settings/hooks/useSettings";
import { useUserInfo } from "@/features/user/hooks/useUserInfo";

const getLocalPath = (
  jobId: number,
  docType: DocumentType,
  company: string,
  title: string
): string => {
  const baseName = `${company.toLowerCase().replace(/ /g, "_")}_${title
    .toLowerCase()
    .replace(/ /g, "_")}_${jobId}`;
  return `${docType}/${baseName}_${docType}.pdf`;
};

const baseFileName = (
  jobId: number,
  docType: DocumentType,
  company: string,
  title: string
): string => {
  const baseName = `${company.toLowerCase().replace(/ /g, "_")}_${title
    .toLowerCase()
    .replace(/ /g, "_")}_${jobId}`;
  return `/${baseName}_${docType}`;
};

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
  const [isApiLoading, setIsApiLoading] = useState(false);
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
      baseFileName(jobId, docType, companyName, jobTitle) + ".pdf";
    const jsonPath =
      baseFileName(jobId, docType, companyName, jobTitle) + ".json";
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
    console.log("Effect triggered", {
      serverStatus,
      jobId,
      downloadUrl: serverStatus?.downloadUrl,
    });
    if (!serverStatus || serverStatus.status !== "COMPLETED") return;
    if (!serverStatus.downloadUrl || !serverStatus.changesUrl) return;
    if (!jobId) return;
    const downloadAndSave = async () => {
      console.log("Downloading and saving file...");
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
          baseFileName(jobId, docType, companyName, jobTitle) + ".pdf";
        const jsonPath =
          baseFileName(jobId, docType, companyName, jobTitle) + ".json";
        await window.appAPI.files.saveFile(pdfPath, arrayBuffer);
        await window.appAPI.files.saveJsonFile(jsonPath, jsonData);
        await checkFile();
      } catch (err) {
        console.error("Failed to download and save file:", err);
        setApiError("Failed to download the new document.");
      }
    };
    downloadAndSave();
    console.log("File downloaded and saved.");
  }, [serverStatus, jobId, docType, companyName, jobTitle, checkFile]);

  const generate = useCallback(async () => {
    if (!jobId || !user || !settings || !userProfile) {
      setApiError("User, Job ID, Settings, or Profile are not loaded.");
      return;
    }

    setIsApiLoading(true);
    setApiError(null);

    try {
      const token = await user.getIdToken();
      const feature =
        docType === "resume" ? "resumeGeneration" : "coverLetterGeneration";
      const llmProvider = settings.featureAssignments[feature];

      const response = await generateDocumentApi(
        docType,
        userProfile,
        jobId,
        llmProvider,
        settings,
        token
      );

      addPendingDocument(String(response.jobId), docType);
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "An unknown error occurred.";
      console.error(
        `Failed to start ${docType} generation for job ${jobId}:`,
        errorMessage
      );
      setApiError(`API Error: ${errorMessage}`);
    } finally {
      setIsApiLoading(false);
    }
  }, [jobId, docType, user, settings, addPendingDocument, userProfile]);

  const displayStatus = useMemo(() => {
    if (isCheckingFile) return "checking";
    if (isApiLoading || serverStatus?.status === "PENDING") return "generating";
    if (serverStatus?.status === "FAILED") return "failed";
    if (fileExists) return "present";
    return "idle";
  }, [isCheckingFile, isApiLoading, serverStatus, fileExists]);

  return {
    displayStatus,
    localPdfPath,
    localJsonData,
    generate,
    error: apiError || serverStatus?.error,
  };
};
