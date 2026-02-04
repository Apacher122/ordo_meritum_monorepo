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
 * Generates a document based on the given document request.
 *
 * This function takes a document request object, the type of the document to be generated,
 * the job ID, the LLM provider, the settings object, and the authentication token.
 * It returns a promise that resolves to a QueueJobResponse object, which contains the job ID
 * and the status of the document generation process.
 *
 * @throws {Error} If the API key for the specified LLM provider is not set.
 * @throws {Error} If the document type is not supported for generation.
 */
export const generateDocument = async (
  docType: DocumentType,
  documentRequest: DocumentRequestBody,
  jobId: number,
  llmProvider: LlmProvider,
  settings: Settings,
  token: string
): Promise<QueueJobResponse> => {
  const featureKey = docType === "resume" ? "resumeGeneration" : "coverLetterGeneration";
  const assignment = settings.featureAssignments[featureKey];
  const providerKeys = settings.apiKeys[llmProvider] || [];
  const apiKey = providerKeys[assignment.keyIndex];

  if (!apiKey) {
    throw new Error(`API key for ${llmProvider} (Index: ${assignment.keyIndex + 1}) is not set.`);
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