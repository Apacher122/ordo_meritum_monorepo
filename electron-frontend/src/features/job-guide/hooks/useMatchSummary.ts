import { useCallback, useEffect, useState } from "react";

import { MatchSummaryResponse } from "@/shared/types/job-guide";
import { getMatchSummary as getMatchSummaryApi } from "../api/getMatchSummary";
import { useApplication } from "@/features/applications/providers/ApplicationProvider";
import { useAuth } from "@/features/auth/providers/AuthProvider";
import { useSettings } from "@/features/settings/hooks/useSettings";

/**
 * Custom hook to manage fetching and caching for a job match summary.
 * It first attempts to load a summary from the local file system. If not found,
 * it provides a function to fetch a new summary from the API and save it locally.
 * @returns {UseMatchSummaryReturn} An object containing the match summary state and management functions.
 */
export const useMatchSummary = () => {
  const { selectedId, selectedJob } = useApplication();
  const { settings } = useSettings();
  const { user } = useAuth();
  const [matchSummary, setMatchSummary] = useState<MatchSummaryResponse | null>(
    null
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasLocalSummary, setHasLocalSummary] = useState(false);

  const getFileName = useCallback(() => {
    if (!selectedJob) return null;
    return `match-summary/${selectedJob.CompanyProperName.toLowerCase().replaceAll(
      " ",
      "_"
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
      const llmProvider = settings.featureAssignments.matchSummary;
      const summary = await getMatchSummaryApi(
        selectedId,
        llmProvider,
        settings,
        token
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
  }, [selectedId, settings, user, getFileName]);

  return { matchSummary, loading, error, getMatchSummary, hasLocalSummary };
};
