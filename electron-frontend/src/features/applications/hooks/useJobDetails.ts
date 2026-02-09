// Location: Gemini/electron-frontend/src/features/applications/hooks/useJobDetails.ts

import { useEffect, useMemo, useState } from "react";

import { useApplication } from "../providers/ApplicationProvider";

/**
 * A hook to fetch job details given a role id.
 *
 * It handles fetching, updating, deleting, and calculating metrics for the applications.
 * @param {roleId: number | string | undefined} The role id of the job to fetch details for.
 * @returns {{ job: AppliedJob | null, isLoading: boolean, error: string | null }}
 */
export const useJobDetails = (id: number | string | undefined) => {
  const { jobs, fetchJobDetails } = useApplication();
  const [isFetching, setIsFetching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const roleId = typeof id === 'string' ? Number.parseInt(id, 10) : id;
  const job = useMemo(() => jobs.find((j) => j.RoleID === roleId), [jobs, roleId]);

  useEffect(() => {
    if (!roleId || Number.isNaN(roleId)) return;

    if (!job?.Description) {
      const loadDetails = async () => {
        setIsFetching(true);
        setError(null);
        try {
          await fetchJobDetails(roleId);
        } catch (err: any) {
          setError(err.message || "Failed to load job details");
        } finally {
          setIsFetching(false);
        }
      };

      loadDetails();
    }
  }, [roleId, job?.Description, fetchJobDetails]);

  return {
    job: job || null,
    isLoading: isFetching,
    error
  };
};