import "@/assets/styles/Components/Layouts/CogitatorView.css";

import { ApplicationStatus, statusOptions } from "../types";
import React, { useEffect, useMemo, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { ApplicationListView } from "../components/ApplicationListView";
import { ApplicationMetrics } from "../components/ApplicationMetrics";
import { SearchHeaderControls } from "../components/SearchHeaderControls";
import { useApplication } from "../providers/ApplicationProvider";
import { useNavigate } from "react-router-dom";

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
  const [searchQuery, setSearchQuery] = useState(""); 
  
  const [showMetrics, setShowMetrics] = useState(false);
  const [statusFilter, setStatusFilter] = useState<ApplicationStatus | "All">("All");
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();
  const navigate = useNavigate();
  
  const headerControls = useMemo(() => (
    <>
      <button 
        onClick={() => navigate("/applicant-data")} 
        className="button primary-button"
      >
        View Applicant Data
      </button>
      
      <SearchHeaderControls onSearch={setSearchQuery} />
      
      <select 
        value={statusFilter} 
        onChange={(e) => setStatusFilter(e.target.value as any)}
        className="status-filter-dropdown"
        aria-label="Filter by status"
      >
        <option value="All">All Statuses</option>
        {statusOptions.map(status => (
          <option key={status} value={status}>{status}</option>
        ))}
      </select>
      <button onClick={() => setShowMetrics(true)} className="button">
        Show Metrics
      </button>      
    </>
  ), [statusFilter, navigate]);
  
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
    let filtered = jobs;
    if (searchQuery) {
      filtered = filtered.filter(
        (job) =>
          job.CompanyName.toLowerCase().includes(searchQuery.toLowerCase()) ||
          job.JobTitle.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    if (statusFilter !== "All") {
      filtered = filtered.filter((job) => job.ApplicationStatus === statusFilter);
    }
    return filtered;
  }, [jobs, searchQuery, statusFilter]);

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