import { LlmProvider } from "@/shared/types";
import { MatchSummaryResponse } from "@/shared/types/job-guide";
import { Settings } from "@/features/settings/types/types";
import { apiRequest } from "@/shared/utils/requests";
import { encryptData } from "@/shared/lib/encryption";

/**
 * Fetches a job match summary from the server for a given job ID.
 * It encrypts the user's API key for the specified LLM provider and includes it
 * in the request header for secure, server-side processing.
 *
 * @param {number} jobId - The unique identifier for the job.
 * @param {LlmProvider} llmProvider - The LLM provider designated to generate the summary.
 * @param {Settings} settings - The application settings object containing API keys.
 * @param {string} token - The user's authentication token.
 * @returns {Promise<MatchSummaryResponse>} A promise that resolves to the match summary response.
 * @throws {Error} If the API key for the specified provider is not set.
 */
export const getMatchSummary = async (
  jobId: number,
  llmProvider: LlmProvider,
  settings: Settings,
  token: string
): Promise<MatchSummaryResponse> => {
  const apiKey = settings.apiKeys[llmProvider];
  if (!apiKey) {
    throw new Error(`API key for ${llmProvider} is not set.`);
  }

  const body = {
    payload: {
      job_id: jobId,
    },
    options: {
      llm: llmProvider.toLowerCase(),
    },
  };
  const encryptedKey = await encryptData(apiKey);
  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
    "X-Encrypted-API-Key": encryptedKey,
    "Content-Type": "application/json",
  };

  return apiRequest<MatchSummaryResponse>(`api/secure/match-summary`, {
    method: "POST",
    headers,
    body: body,
  });
};
