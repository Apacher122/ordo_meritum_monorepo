import React, { useEffect, useMemo, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { ApplicationListView } from "../components/ApplicationListView";
import { ApplicationMetrics } from "../components/ApplicationMetrics";
import { SearchHeaderControls } from "../components/SearchHeaderControls";
import { useApplication } from "../providers/ApplicationProvider";

export const ApplicationListPage: React.FC = () => {
  const { jobs, metrics, loading, error, updateJobStatus, updateJobDate, removeJob } = useApplication();
  const [searchQuery, setSearchQuery] = useState("");
  const [showMetrics, setShowMetrics] = useState(false);
  
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();
  
  
  const headerControls = useMemo(() => (
    <>
      <SearchHeaderControls onSearch={setSearchQuery} initialQuery={searchQuery} />
      <button onClick={() => setShowMetrics(true)} className="button">
        Show Metrics
      </button>      
    </>
  ), [searchQuery]); 
  
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
    <div className="application-list-page">
      <ApplicationListView
        jobs={filteredJobs}
        onStatusUpdate={updateJobStatus}
        onDateUpdate={updateJobDate}
        onDelete={removeJob}
      />

      {showMetrics && <ApplicationMetrics metrics={metrics} onClose={() => setShowMetrics(false)} />}
    </div>
  );
};

