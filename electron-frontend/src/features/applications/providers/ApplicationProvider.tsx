import { ApplicationStatus, AppliedJob } from "../types";
import React, {
  ReactNode,
  createContext,
  useContext,
  useMemo,
} from "react";

import { ApplicationMetricsData } from "../components/ApplicationMetrics";
import { useApplicationList } from "../hooks/useApplicationList";

interface ApplicationContextType {
  jobs: AppliedJob[];
  selectedId: number | null;
  setSelectedId: (id: number | null) => void;
  selectedJob: AppliedJob | null;
  updateJobStatus: (roleId: number, newStatus: ApplicationStatus) => void;
  updateJobDate: (roleId: number, newDate: Date) => void;
  removeJob: (roleId: number) => void;
  metrics: ApplicationMetricsData;   loading: boolean;
  error: string | null;
}
const ApplicationContext = createContext<ApplicationContextType | undefined>(undefined);


export const ApplicationProvider = ({ children }: { children: ReactNode }) => {
  const {
    jobs,
    selectedId,
    setSelectedId,
    updateJobStatus,
    updateJobDate,
    removeJob,
    metrics,
    loading,
    error,
  } = useApplicationList();

  const selectedJob = useMemo(() => {
    return jobs.find((job) => job.RoleID === selectedId) || null;
  }, [jobs, selectedId]);

  const value = {
    jobs,
    selectedId,
    setSelectedId,
    selectedJob,
    updateJobStatus,
    updateJobDate,
    removeJob,
    metrics,
    loading,
    error,
  };

  return (
    <ApplicationContext.Provider value={value}>
      {children}
    </ApplicationContext.Provider>
  );
};

export const useApplication = () => {
  const context = useContext(ApplicationContext);
  if (!context) {
    throw new Error("useApplication must be used within an ApplicationProvider");
  }
  return context;
};

