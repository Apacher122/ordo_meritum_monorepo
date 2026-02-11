import "@/assets/styles/pages/InfoPage.css";

import { ApplicationData, ApplicationStatus, statusOptions } from "../types";
import React, { useEffect, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { useApplication } from "../providers/ApplicationProvider";
import { useJobDetails } from "../hooks/useJobDetails";
import { useParams } from "react-router-dom";

/**
 * A page to display the details of a job application.
 * @param {ApplicationData} job - The job application to display.
 * @param {boolean} isLoading - Whether the job application is currently being loaded.
 * @param {string | null} error - Any error that occurred while loading the job application.
 * @param {ApplicationData} metrics - The metrics of the job application.
 * @param {(e: React.ChangeEvent<HTMLInputElement>) => void} handleMetricChange - A function to handle changes to the metrics of the job application.
 * @param {() => void} saveMetrics - A function to save the metrics of the job application locally.
 */
export const JobDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { updateJobStatus } = useApplication();
  
  const { job, isLoading, error } = useJobDetails(id);

  const [metrics, setMetrics] = useState<ApplicationData>({
    hasSummary: false,
    hasCoverLetter: false,
    isVeteran: false,
    isDisabled: false,
    raceMentioned: false
  });

  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  useEffect(() => {
    const saved = localStorage.getItem(`metrics-${id}`);
    if (saved) setMetrics(JSON.parse(saved));
  }, [id]);

  useEffect(() => {
    if (job) {
      setHeaderTitle(job.CompanyProperName);
      setHeaderSubtitle(job.JobTitle);
      setHeaderControls(
        <button onClick={() => updateJobStatus(job.RoleID, "Open")} className="button success-button">
          Applied
        </button>
      );
    }
  }, [job, setHeaderTitle, setHeaderSubtitle, setHeaderControls, updateJobStatus]);

  const handleMetricChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, checked } = e.target;
    setMetrics(prev => ({ ...prev, [name]: checked }));
  };

  const saveMetrics = () => {
    localStorage.setItem(`metrics-${id}`, JSON.stringify(metrics));
    alert("Metrics saved locally.");
  };
  
  if (isLoading) return <div style={{ padding: '20px' }}>Loading job details...</div>;
  if (error) return <div style={{ padding: '20px', color: 'red' }}>Error: {error}</div>;
  if (!job) return <div style={{ padding: '20px' }}>Job application not found.</div>;

  return (
    <div className="job-details-page" style={{ padding: '20px', color: 'var(--text-main)' }}>
      <section className="details-section">
        <div className="info-grid">
          <p><strong>Salary Range:</strong> {job.SalaryRange || "N/A"}</p>
          <p><strong>Status:</strong> 
            <select 
              value={job.ApplicationStatus} 
              onChange={(e) => updateJobStatus(job.RoleID, e.target.value as ApplicationStatus)}
              className="status-dropdown"
            >
              {statusOptions.map(opt => <option key={opt} value={opt}>{opt}</option>)}
            </select>
          </p>
          {job.Website && <a href={job.Website} target="_blank" rel="noreferrer" className="link">Link to Job Post</a>}
        </div>
        <div className="description-box" style={{ marginTop: '15px' }}>
          <h4>Job Description</h4>
          <p style={{ whiteSpace: 'pre-wrap' }}>{job.Description}</p>
        </div>
      </section>

      <hr />

      <section className="details-section">
        <div className="form-grid">
          <div className="column">
            <p><strong>Years of Experience:</strong> {job.YearsOfExp || "N/A"}</p>
            <p><strong>Education:</strong> {job.EducationLevel || "N/A"}</p>
            <p><strong>Applicants:</strong> {job.ApplicantCount || "N/A"}</p>
            <p><strong>Posted:</strong> {job.PostAge || "N/A"}</p>
            <p><strong>Skills:</strong> {job.Requirements?.join(", ") || "N/A"}</p>
            <p><strong>Nice to Haves:</strong> {job.NiceToHaves?.join(", ") || "N/A"}</p>
            <p><strong>Tools/Tech:</strong> {job.Tools?.join(", ") || "N/A"}</p>
            <p><strong>Languages:</strong> {job.ProgrammingLanguages?.join(", ") || "N/A"}</p>
          </div>
          <div className="column">
            <p><strong>Frameworks:</strong> {job.FrameworksAndLibraries?.join(", ") || "N/A"}</p>
            <p><strong>Databases:</strong> {job.Databases?.join(", ") || "N/A"}</p>
            <p><strong>Cloud:</strong> {job.CloudTechnologies?.join(", ") || "N/A"}</p>
            <p><strong>Keywords:</strong> {job.IndustryKeywords?.join(", ") || "N/A"}</p>
            <p><strong>Soft Skills:</strong> {job.SoftSkills?.join(", ") || "N/A"}</p>
            <p><strong>Certs:</strong> {job.Certifications?.join(", ") || "N/A"}</p>
            <p><strong>Culture:</strong> {job.CompanyCulture || "N/A"}</p>
            <p><strong>Values:</strong> {job.CompanyValues || "N/A"}</p>
          </div>
        </div>
      </section>

      <hr />

      <section className="details-section">
        <h4>Application Details</h4>
        <div className="form-grid">
          <div className="column">
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="hasSummary" checked={metrics.hasSummary} onChange={handleMetricChange} />
              <span>Resume includes summary</span>
            </label>
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="hasCoverLetter" checked={metrics.hasCoverLetter} onChange={handleMetricChange} />
              <span>Cover Letter sent</span>
            </label>
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="isVeteran" checked={metrics.isVeteran} onChange={handleMetricChange} />
              <span>Protected Veteran</span>
            </label>
          </div>
          <div className="column">
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="isDisabled" checked={metrics.isDisabled} onChange={handleMetricChange} />
              <span>Person with Disability</span>
            </label>
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="raceMentioned" checked={metrics.raceMentioned} onChange={handleMetricChange} />
              <span>Race Mentioned</span>
            </label>
          </div>
        </div>
        <button onClick={saveMetrics} className="button primary-button" style={{ marginTop: '20px' }}>
          Save Metrics
        </button>
      </section>
    </div>
  );
};