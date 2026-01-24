import { ApplicationStatus, AppliedJob } from "../types";
import React, {
  ReactNode,
  createContext,
  useContext,
  useMemo,
} from "react";

import { useApplicationList } from "../hooks/useApplicationList";

interface ApplicationContextType {
jobs: AppliedJob[];
  selectedId: number | null;
  setSelectedId: (id: number | null) => void;
  selectedJob: AppliedJob | null; // <-- This is correctly KEPT
  updateJobStatus: (roleId: number, newStatus: ApplicationStatus) => void;
  updateJobDate: (roleId: number, newDate: Date) => void;
  removeJob: (roleId: number) => void;
  loading: boolean;
  error: string | null;
}
const ApplicationContext = createContext<ApplicationContextType | undefined>(undefined);

/**
 * Provides application-related state and actions to its children components.
 * This context centralizes the management of job application data.
 * @param {object} props - The component props.
 * @param {ReactNode} props.children - The child components to render.
 * @returns {JSX.Element}
 */
export const ApplicationProvider = ({ children }: { children: ReactNode }) => {
  const {
    jobs,
    selectedId,
    setSelectedId,
    updateJobStatus,
    updateJobDate,
    removeJob,
    loading,
    error,
  } = useApplicationList();

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
  }), [
    jobs,
    selectedId,
    setSelectedId,
    selectedJob,
    updateJobStatus,
    updateJobDate,
    removeJob,
    loading,
    error
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

