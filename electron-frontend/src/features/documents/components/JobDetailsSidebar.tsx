import "@/assets/styles/Components/JobDetailsSidebar.css";

import { ApplicationStatus, AppliedJob, statusOptions } from "@/features/applications/types";

import React from "react";
import { useApplication } from "@/features/applications/providers/ApplicationProvider";
import { useJobDetails } from "@/features/applications/hooks/useJobDetails";

interface JobDetailsSidebarProps {
  roleId: number;
}

export const JobDetailsSidebar: React.FC<JobDetailsSidebarProps> = ({ roleId }) => {
  const { updateJobStatus } = useApplication();
  const { job, isLoading, error } = useJobDetails(roleId);

  return (
    <div className="job-details-sidebar">
      <div className="sidebar-header">
        <h3>Job Details</h3>
        <div className="status-control">
          <label>Status:</label>
          <select
            value={job?.ApplicationStatus}
            onChange={(e) => updateJobStatus(job?.RoleID as number, e.target.value as ApplicationStatus)}
            className="status-dropdown-small"
          >
            {statusOptions.map((opt) => (
              <option key={opt} value={opt}>
                {opt}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="sidebar-content">
        <section className="sidebar-section">
          <h4>Overview</h4>
          <p><strong>Company:</strong> {job?.CompanyProperName}</p>
          <p><strong>Role:</strong> {job?.JobTitle}</p>
          <p><strong>Salary:</strong> {job?.SalaryRange || "N/A"}</p>
          {job?.Website && (
            <p>
              <a href={job.Website} target="_blank" rel="noreferrer" className="link">
                View Original Post
              </a>
            </p>
          )}
        </section>

        <section className="sidebar-section">
          <h4>Skills & Requirements</h4>
          {job?.Requirements && job?.Requirements.length > 0 && (
            <div className="tag-group">
              <strong>Requirements:</strong>
              <ul>
                {job?.Requirements.map((req, i) => <li key={i}>{req}</li>)}
              </ul>
            </div>
          )}
          
          {job?.NiceToHaves && job?.NiceToHaves.length > 0 && (
            <div className="tag-group">
              <strong>Nice to Have:</strong>
              <ul>
                {job.NiceToHaves.map((item, i) => <li key={i}>{item}</li>)}
              </ul>
            </div>
          )}

          {job?.Tools && job?.Tools.length > 0 && (
             <div className="tag-group">
               <strong>Tools:</strong> <span className="tag-list">{job.Tools.join(", ")}</span>
             </div>
          )}
        </section>

        
        <section className="sidebar-section">
          <h4>Description</h4>
          <p className="description-text">{job?.Description}</p>
        </section>
      </div>
    </div>
  );
};