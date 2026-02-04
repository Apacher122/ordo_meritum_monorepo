import { LlmProvider } from "@/shared/types";
import { MatchSummaryResponse } from "@/shared/types/job-guide";
import { Settings } from "@/features/settings/types/types";
import { apiRequest } from "@/shared/utils/requests";
import { encryptData } from "@/shared/lib/encryption";

/**
 * Retrieves a match summary for a job.
 *
 * @param {number} jobId - The ID of the job to retrieve the match summary for.
 * @param {LlmProvider} llmProvider - The LLM provider to use for generating the match summary.
 * @param {Settings} settings - The application settings, containing the API keys and assignments.
 * @param {string} token - The authentication token to use for the request.
 * @returns {Promise<MatchSummaryResponse>} A promise resolving to the match summary response.
 * @throws {Error} If the API key for the specified LLM provider is not set.
 */
export const getMatchSummary = async (
  jobId: number,
  llmProvider: LlmProvider,
  settings: Settings,
  token: string,
): Promise<MatchSummaryResponse> => {
  const assignment = settings.featureAssignments.matchSummary;
  const providerKeys = settings.apiKeys[llmProvider] || [];
  const apiKey = providerKeys[assignment.keyIndex];

  if (!apiKey) {
    throw new Error(
      `API key for ${llmProvider} (Index: ${assignment.keyIndex + 1}) is not set.`,
    );
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
