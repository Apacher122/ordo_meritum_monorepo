import { LlmProvider } from '@/shared/types';
import { DocumentType } from '../types';
import { apiRequest } from '@/shared/utils/requests';
import { encryptData } from '@/shared/lib/encryption';
import { Settings } from '@/features/settings/types/types';
import { DocumentRequestBody } from '../types';

export interface QueueJobResponse {
  jobId: number;
  status: string;
}

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
    "Authorization": `Bearer ${token}`,
    "X-Encrypted-API-Key": encryptedKey,
    'Content-Type': 'application/json',
  };

  let payload: any;
  if (docType === "resume") {
    payload = {
      userInfo: documentRequest.userInfo,
      educationInfo: documentRequest.education,
      resume: documentRequest.resume,
      additionalInfo: documentRequest.aboutMe,
    };
  } else if (docType === "cover-letter") {
    payload = { ...documentRequest.coverLetter, userInfo: documentRequest.userInfo };
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
