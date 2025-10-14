import * as api from "../api";

import { ApplicationStatus, AppliedJob } from "../types";
import { denormalizeStatus, normalizeStatus } from "../utils/statusMappings";
import { useCallback, useEffect, useMemo, useState } from "react";

import { ApplicationMetricsData } from "../components/ApplicationMetrics";

export const useApplicationList = () => {
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
      ApplicationStatus: normalizeStatus(job.ApplicationStatus as string),
      InitialApplicationDate: correctedDate,
    };
  };

  const metrics = useMemo((): ApplicationMetricsData => {
    const appsSent = jobs.filter((j) => j.ApplicationStatus !== "Not Applied");
    const applicationsSent = appsSent.length;
    if (applicationsSent === 0) {
      return {
        applicationsSent: 0,
        rejections: 0,
        ghosts: 0,
        stillOpen: 0,
        interviewing: 0,
        offers: 0,
        interviewRate: 0,
        rejectionRate: 0,
        ghostedRate: 0,
        openAppsRate: 0,
        appsSentToday: 0,
        avgAppsPerDay: "0.0",
      };
    }

    const rejections = appsSent.filter(
      (j) => j.ApplicationStatus === "Rejected"
    ).length;
    const ghosts = appsSent.filter(
      (j) => j.ApplicationStatus === "Ghosted"
    ).length;
    const stillOpen = appsSent.filter(
      (j) =>
        j.ApplicationStatus === "Open" || j.ApplicationStatus === "Interviewing"
    ).length;
    const interviewing = appsSent.filter(
      (j) => j.ApplicationStatus === "Interviewing"
    ).length;
    const offers = appsSent.filter(
      (j) => j.ApplicationStatus === "Offered"
    ).length;

    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const appsSentToday = appsSent.filter((j) => {
      const appDate = new Date(j.InitialApplicationDate);
      appDate.setHours(0, 0, 0, 0);
      return appDate.getTime() === today.getTime();
    }).length;

    const firstAppDate = appsSent.reduce((oldest, job) => {
      const jobDate = new Date(job.InitialApplicationDate);
      return jobDate < oldest ? jobDate : oldest;
    }, new Date());
    const daysSinceFirstApp = Math.max(
      1,
      Math.ceil(
        (new Date().getTime() - firstAppDate.getTime()) / (1000 * 60 * 60 * 24)
      )
    );
    const avgAppsPerDay = (applicationsSent / daysSinceFirstApp).toFixed(1);

    return {
      applicationsSent,
      rejections,
      ghosts,
      stillOpen,
      interviewing,
      offers,
      interviewRate: (interviewing / applicationsSent) * 100,
      rejectionRate: (rejections / applicationsSent) * 100,
      ghostedRate: (ghosts / applicationsSent) * 100,
      openAppsRate: (stillOpen / applicationsSent) * 100,
      appsSentToday,
      avgAppsPerDay,
    };
  }, [jobs]);

  const fetchJobs = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getApplications();
      const sortedData = data.sort(
        (a, b) =>
          new Date(b.InitialApplicationDate).getTime() -
          new Date(a.InitialApplicationDate).getTime()
      );
      setJobs(sortedData.map(transformJobData));
    } catch (err) {
      setError("Failed to load applications.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchJobs();
  }, [fetchJobs]);

  const updateJobStatus = useCallback(
    async (roleId: number, newStatus: ApplicationStatus) => {
      const originalJobs = jobs;
      const backendStatus = denormalizeStatus(newStatus);
      setJobs((prev) =>
        prev.map((j) =>
          j.RoleID === roleId ? { ...j, ApplicationStatus: newStatus } : j
        )
      );
      try {
        await api.updateApplication(roleId, {
          status: backendStatus,
          date: null,
        });
      } catch (err) {
        setJobs(originalJobs);
        setError("Failed to update status.");
      }
    },
    [jobs]
  );

  const updateJobDate = useCallback(
    async (roleId: number, newDate: Date) => {
      const originalJobs = jobs;
      setJobs((prev) =>
        prev.map((j) =>
          j.RoleID === roleId ? { ...j, InitialApplicationDate: newDate } : j
        )
      );
      try {
        await api.updateApplication(roleId, { status: null, date: newDate });
      } catch (err) {
        setJobs(originalJobs);
        setError("Failed to update date.");
      }
    },
    [jobs]
  );

  const removeJob = useCallback(
    async (roleId: number) => {
      const originalJobs = jobs;
      setJobs((prev) => prev.filter((j) => j.RoleID !== roleId));
      try {
        await api.deleteApplication(roleId);
      } catch (err) {
        setJobs(originalJobs);
        setError("Failed to delete application.");
      }
    },
    [jobs]
  );

  return {
    jobs,
    metrics,
    loading,
    error,
    selectedId,
    setSelectedId,
    updateJobStatus,
    updateJobDate,
    removeJob,
  };
};
