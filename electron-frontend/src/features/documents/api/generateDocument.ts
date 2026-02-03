import { DocumentRequestBody, DocumentType } from "../types";

import { LlmProvider } from "@/shared/types";
import { Settings } from "@/features/settings/types/types";
import { apiRequest } from "@/shared/utils/requests";
import { encryptData } from "@/shared/lib/encryption";

export interface QueueJobResponse {
  jobId: number;
  status: string;
}

/**
 * Generates a document of the specified type using an LLM.
 *
 * @param {DocumentType} docType The type of document to generate.
 * @param {DocumentRequestBody} documentRequest The request body containing user information, resume content, and other relevant data.
 * @param {number} jobId The job ID for the generated document.
 * @param {LlmProvider} llmProvider The LLM provider to use for generation.
 * @param {Settings} settings The application settings containing API keys for LLM providers.
 * @param {string} token The Bearer token for API requests.
 *
 * @returns {Promise<QueueJobResponse>} A promise resolving to a QueueJobResponse containing the job ID and status of the generated document.
 */
export const generateDocument = async (
  docType: DocumentType,
  documentRequest: DocumentRequestBody,
  jobId: number,
  llmProvider: LlmProvider,
  settings: Settings,
  token: string
): Promise<QueueJobResponse> => {
  const apiKey = settings.apiKeys[llmProvider];
  if (!apiKey) {
    throw new Error(`API key for ${llmProvider} is not set.`);
  }

  const encryptedKey = await encryptData(apiKey);

  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
    "X-Encrypted-API-Key": encryptedKey,
    "Content-Type": "application/json",
  };

  const sanitizedResume = {
    ...documentRequest.resume,
    experiences: documentRequest.resume.experiences?.map((exp) => ({
      ...exp,
      bulletPoints: exp.bulletPoints?.map((bp) => bp.text),
    })),
    projects: documentRequest.resume.projects?.map((proj) => ({
      ...proj,
      bulletPoints: proj.bulletPoints?.map((bp) => bp.text),
    })),
  };

  let payload: any;
  if (docType === "resume") {
    payload = {
      userInfo: documentRequest.userInfo,
      educationInfo: documentRequest.education,
      resume: sanitizedResume,
      additionalInfo: documentRequest.aboutMe,
      writingSamples: documentRequest.writingSamples,
    };
  } else if (docType === "cover-letter") {
    payload = {
      ...documentRequest.coverLetter,
      userInfo: documentRequest.userInfo,
      educationInfo: documentRequest.education,
      additionalInfo: documentRequest.aboutMe,
      writingSamples: documentRequest.writingSamples,
    };
  } else {
    throw new Error("Unsupported document type for generation.");
  }

  const body = {
    payload: payload,
    options: {
      jobId: jobId,
      docType: docType,
      llm: llmProvider.toLowerCase(),
    },
  };

  return await apiRequest<QueueJobResponse>(`api/secure/documents/${docType}`, {
    method: "POST",
    headers,
    body: body,
  });
};
