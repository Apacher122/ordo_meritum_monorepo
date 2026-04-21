import { useCallback, useEffect, useState } from "react";

import { MatchSummaryResponse } from "@/shared/types/job-guide";
import { getMatchSummary as getMatchSummaryApi } from "../api/getMatchSummary";
import { useApplication } from "@/features/applications/providers/ApplicationProvider";
import { useAuth } from "@/features/auth/providers/AuthProvider";
import { useSettings } from "@/features/settings/hooks/useSettings";

/**
 * A hook that retrieves a match summary for a job.
 *
 * It first checks if a user is authenticated and if not, it returns an error.
 * Then, it attempts to retrieve the match summary from the local file system and if that fails, it returns an error.
 * If the match summary is successfully retrieved, it returns the match summary, whether it's loading, and whether there's an error.
 *
 * @returns {{ matchSummary: MatchSummaryResponse | null, loading: boolean, error: string | null, getMatchSummary: (jobId: number) => Promise<void>, hasLocalSummary: boolean }}
 */
export const useMatchSummary = () => {
  const { selectedId, selectedJob } = useApplication();
  const { settings, getValidApiKey } = useSettings();
  const { user } = useAuth();
  const [matchSummary, setMatchSummary] = useState<MatchSummaryResponse | null>(
    null,
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasLocalSummary, setHasLocalSummary] = useState(false);

  const getFileName = useCallback(() => {
    if (!selectedJob) return null;
    return `match-summary/${selectedJob.CompanyProperName.toLowerCase().replaceAll(
      " ",
      "_",
    )}_${selectedJob.JobTitle.toLowerCase().replaceAll(" ", "_")}_${
      selectedJob.RoleID
    }_match_summary.json`;
  }, [selectedJob]);

  const loadLocalSummary = useCallback(async () => {
    const fileName = getFileName();
    if (!fileName) {
      setHasLocalSummary(false);
      return;
    }

    setLoading(true);
    const result = await window.appAPI.files.readJsonFile(fileName);
    if (result.success && result.data) {
      setMatchSummary(result.data);
      setHasLocalSummary(true);
    } else {
      setMatchSummary(null);
      setHasLocalSummary(false);
    }
    setLoading(false);
  }, [getFileName]);

  useEffect(() => {
    loadLocalSummary();
  }, [loadLocalSummary]);

  const getMatchSummary = useCallback(async () => {
    if (!selectedId || !settings || !user) {
      setMatchSummary(null);
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const token = await user.getIdToken();
      const assignment = settings.featureAssignments.matchSummary;
      const llmProvider = assignment.provider;

      if (llmProvider === "None") {
        throw new Error(
          "No LLM provider assigned for Match Summary in settings.",
        );
      }

      const apiKey = await getValidApiKey(llmProvider);

      if (!apiKey) {
        throw new Error(
          `Rate limit reached or no keys available for ${llmProvider}.`,
        );
      }

      const summary = await getMatchSummaryApi(
        selectedId,
        llmProvider,
        apiKey,
        token,
      );
      setMatchSummary(summary);
      const fileName = getFileName();
      if (fileName) {
        await window.appAPI.files.saveJsonFile(fileName, summary);
      }
      setHasLocalSummary(true);
    } catch (err: any) {
      setError(err.message || "An unknown error occurred.");
    } finally {
      setLoading(false);
    }
  }, [selectedId, settings, user, getFileName, getValidApiKey]);

  return { matchSummary, loading, error, getMatchSummary, hasLocalSummary };
};
