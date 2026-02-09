import { ApplicationStatus, AppliedJob } from "../types";
import React, {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";

import { getApplicationById } from "../api/getApplications";
import { useApplicationList } from "../hooks/useApplicationList";

interface ApplicationContextType {
  jobs: AppliedJob[];
  selectedId: number | null;
  setSelectedId: (id: number | null) => void;
  selectedJob: AppliedJob | null;
  updateJobStatus: (roleId: number, newStatus: ApplicationStatus) => void;
  updateJobDate: (roleId: number, newDate: Date) => void;
  removeJob: (roleId: number) => void;
  loading: boolean;
  error: string | null;
  fetchJobDetails: (id: number) => Promise<void>;
}
const ApplicationContext = createContext<ApplicationContextType | undefined>(undefined);

/**
 * A React Context Provider that manages the state and operations for the user's list of job applications.
 * @param {children: ReactNode} The children elements to be rendered within the context provider.
 * @returns {React.Context<ApplicationContextType>} A React context containing the application list state and management functions.
 */
export const ApplicationProvider = ({ children }: { children: ReactNode }) => {
  const {
    jobs: rawJobs,
    selectedId,
    setSelectedId,
    updateJobStatus,
    updateJobDate,
    removeJob,
    loading,
    error,
  } = useApplicationList();

  const [detailedCache, setDetailedCache] = useState<Record<number, Partial<AppliedJob>>>({});

  const jobs = useMemo(() => {
    return rawJobs.map((job) => ({
      ...job,
      ...detailedCache[job.RoleID],
    }));
  }, [rawJobs, detailedCache]);

  const fetchJobDetails = useCallback(async (id: number) => {
    const jobInState = jobs.find(j => j.RoleID === id);
    
    if (jobInState?.Description) {
      return;
    }

    try {
      const response = await getApplicationById(id);
      
      setDetailedCache((prev) => ({
        ...prev,
        [id]: response,
      }));
    } catch (err: any) {
      console.error("Error fetching job details:", err);
      throw err;
    }
  }, [jobs]);

  const selectedJob = useMemo(() => {
    return jobs.find((job) => job.RoleID === selectedId) || null;
  }, [jobs, selectedId]);

  const value = useMemo(() => ({
    jobs,
    selectedId,
    setSelectedId,
    selectedJob,
    updateJobStatus,
    updateJobDate,
    removeJob,
    loading,
    error,
    fetchJobDetails,
  }), [
    jobs,
    selectedId,
    setSelectedId,
    selectedJob,
    updateJobStatus,
    updateJobDate,
    removeJob,
    loading,
    error,
    fetchJobDetails
  ]);

  return (
    <ApplicationContext.Provider value={value}>
      {children}
    </ApplicationContext.Provider>
  );
};

/**
 * Custom hook to access the ApplicationContext.
 * @returns {ApplicationContextType} The application context.
 * @throws {Error} If used outside of an `ApplicationProvider`.
 */
export const useApplication = () => {
  const context = useContext(ApplicationContext);
  if (!context) {
    throw new Error("useApplication must be used within an ApplicationProvider");
  }
  return context;
};

