import { LlmProvider } from "@/shared/types";
import { MatchSummaryResponse } from "@/shared/types/job-guide";
import { Settings } from "@/features/settings/types/types";
import { apiRequest } from "@/shared/utils/requests";
import { encryptData } from "@/shared/lib/encryption";

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
