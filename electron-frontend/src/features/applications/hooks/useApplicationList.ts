import * as api from "../api";

import { ApplicationStatus, AppliedJob } from "../types";
import { denormalizeStatus, normalizeStatus } from "../utils/statusMappings";
import { useCallback, useEffect, useState } from "react";

import { useAuth } from "@/app/appProviders";

/**
 * Custom hook to manage the state and operations for the user's list of job applications.
 * It handles fetching, updating, deleting, and calculating metrics for the applications.
 * @returns {UseApplicationListReturn} An object containing the application list state and management functions.
 */
export const useApplicationList = () => {
  const { user } = useAuth();
  const [jobs, setJobs] = useState<AppliedJob[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const transformJobData = (job: AppliedJob): AppliedJob => {
    const dateFromBackend = new Date(job.InitialApplicationDate);

    const timezoneOffset = dateFromBackend.getTimezoneOffset() * 60000;
    const correctedDate = new Date(dateFromBackend.getTime() + timezoneOffset);
    return {
      ...job,
      ApplicationStatus: normalizeStatus(job.ApplicationStatus),
      InitialApplicationDate: correctedDate,
    };
  };

  const fetchJobs = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getApplications(); //
      const sortedData = [...data].sort(
        (a, b) =>
          new Date(b.InitialApplicationDate).getTime() -
          new Date(a.InitialApplicationDate).getTime()
      );
      const transformedJobs = sortedData.map(transformJobData);
      setJobs(transformedJobs);

      setTimeout(() => {
        localStorage.setItem("jobs", JSON.stringify(transformedJobs));
      }, 0);

      setError(null);
    } catch (err) {
      console.log(err);
      setError("Server offline. Loading applications from local cache.");
      try {
        const cachedJobs = localStorage.getItem("jobs");
        if (cachedJobs) {
          setJobs(JSON.parse(cachedJobs));
        }
      } catch (cacheError) {
        console.log(cacheError);
        setError("Could not load applications from local cache.");
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchJobs();
  }, [fetchJobs]);

  const updateJobStatus = useCallback(
    async (roleId: number, newStatus: ApplicationStatus) => {
      const backendStatus = denormalizeStatus(newStatus);
      if (!user) {
        setError("User, Job ID, Settings, or Profile are not loaded.");
        return;
      }
      const token = await user.getIdToken();
      setJobs((prev) =>
        prev.map((j) =>
          j.RoleID === roleId ? { ...j, ApplicationStatus: newStatus } : j
        )
      );
      try {
        await api.updateApplication(token, {
          job_id: roleId,
          job_title: null,
          website: null,
          application_status: backendStatus,
          interview_count: null,
          initial_application_date: null,
        });
      } catch (err) {
        setError("Failed to update status." + err);
      }
    },
    [user]
  );

  const updateJobDate = useCallback(
    async (roleId: number, newDate: Date) => {
      if (!user) {
        setError("User, Job ID, Settings, or Profile are not loaded.");
        return;
      }
      const token = await user.getIdToken();
      setJobs((prev) =>
        prev.map((j) =>
          j.RoleID === roleId ? { ...j, InitialApplicationDate: newDate } : j
        )
      );
      try {
        await api.updateApplication(token, {
          job_id: roleId,
          job_title: null,
          website: null,
          application_status: null,
          interview_count: null,
          initial_application_date: newDate,
        });
      } catch (err) {
        setError("Failed to update date." + err);
      }
    },
    [user]
  );

  const removeJob = useCallback(
    async (roleId: number) => {
      if (!user) {
        setError("User, Job ID, Settings, or Profile are not loaded.");
        return;
      }
      const token = await user.getIdToken();
      setJobs((prev) => prev.filter((j) => j.RoleID !== roleId));
      try {
        await api.deleteApplication(token, roleId);
      } catch (err) {
        setError("Failed to delete application." + err);
      }
    },
    [user]
  );

  return {
    jobs,
    loading,
    error,
    selectedId,
    setSelectedId,
    updateJobStatus,
    updateJobDate,
    removeJob,
  };
};
