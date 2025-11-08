import "@/assets/styles/Components/Layouts/CogitatorView.css";

import React, { useEffect, useMemo, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { ApplicationListView } from "../components/ApplicationListView";
import { ApplicationMetrics } from "../components/ApplicationMetrics";
import { SearchHeaderControls } from "../components/SearchHeaderControls";
import { useApplication } from "../providers/ApplicationProvider";
import { useDebounce } from "../hooks/useDebounce";

/**
 * Renders the main page for viewing and managing tracked job applications.
 *
 * This component orchestrates the application tracking view. It:
 * - Fetches application data (jobs, loading state, etc.) via the `useApplication` hook.
 * - Manages header content, setting the title and adding controls for search
 * and toggling the metrics modal.
 * - Implements a debounced search to filter the applications list by company
 * or job title without lagging on user input.
 * - Passes the filtered list to `ApplicationListView` for rendering.
 * - Conditionally renders `ApplicationMetrics` for an on-demand metrics overview,
 * passing the complete (unfiltered) jobs list for calculation.
 *
 * @returns {React.ReactElement} The rendered ApplicationListPage component.
 */
export const ApplicationListPage: React.FC = () => {
  const { jobs, loading, error, updateJobStatus, updateJobDate, removeJob } = useApplication();
  const [inputValue, setInputValue] = useState("");
  const searchQuery = useDebounce(inputValue, 300);
  const [showMetrics, setShowMetrics] = useState(false);
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();
  
  const headerControls = useMemo(() => (
    <>
      <SearchHeaderControls 
        value={inputValue} 
        onSearch={setInputValue} 
      />
      <button onClick={() => setShowMetrics(true)} className="button">
        Show Metrics
      </button>      
    </>
  ), [inputValue]);
  
  useEffect(() => {
    setHeaderTitle("My Applications");
    setHeaderSubtitle("Track your job search progress");
    setHeaderControls(headerControls);
    
    return () => {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
      setHeaderControls(null);
    };
  }, [setHeaderTitle, setHeaderSubtitle, setHeaderControls, headerControls]); 

  const filteredJobs = useMemo(() => {
    if (!searchQuery) return jobs;
    return jobs.filter(
      (job) =>
        job.CompanyName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        job.JobTitle.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [jobs, searchQuery]);

  if (loading) return <div>Loading applications...</div>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div className="application-list-page cogitator-view">
      <ApplicationListView
        jobs={filteredJobs}
        onStatusUpdate={updateJobStatus}
        onDateUpdate={updateJobDate}
        onDelete={removeJob}
      />
      {showMetrics && (
        <ApplicationMetrics 
          jobs={jobs} 
          onClose={() => setShowMetrics(false)} 
        />
      )}    
    </div>
  );
};