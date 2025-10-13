import { LlmProvider } from "@/shared/types";
import { apiRequest } from "@/shared/utils/requests";

export interface NewJobPayload {
  companyName: string;
  positionTitle: string;
  url: string;
  description: string;
}

export interface JobSubmissionOptions {
  llmProvider: LlmProvider;
  encryptedApiKey?: string;
}

/**
 * Sends a new job to the backend for analysis.
 * @param {NewJobPayload} jobData - The data of the job to be sent.
 * @param {JobSubmissionOptions} options - Optional parameters to customize the job submission.
 * @returns {Promise<string>} - A promise that resolves to a success message if the job is submitted successfully.
 * @throws {Error} - If the job submission fails.
 */
export const sendJobInfo = async (
  jobData: NewJobPayload,
  options: JobSubmissionOptions
): Promise<string> => {
  const headers: Record<string, string> = {};
  if (options.llmProvider) headers["X-LLM-Provider"] = options.llmProvider;
  if (options.encryptedApiKey) {
    headers["X-Encrypted-API-Key"] = options.encryptedApiKey;
  }

  const response = await apiRequest<{ success: boolean; message: string }>(
    "api/secure/user/send-job-info",
    {
      method: "POST",
      body: jobData,
      headers,
    }
  );

  if (!response.success) {
    throw new Error(`Error sending job info: ${response.message}`);
  }

  return response.message;
};