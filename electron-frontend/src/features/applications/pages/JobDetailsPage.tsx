import "@/assets/styles/pages/InfoPage.css";

import { ApplicationData, ApplicationStatus, statusOptions } from "../types";
import React, { useEffect, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { useApplication } from "../providers/ApplicationProvider";
import { useJobDetails } from "../hooks/useJobDetails";
import { useParams } from "react-router-dom";

/**
* A page to display the details of a job application.
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
    raceMentioned: false,
    isEasyApply: false,
    timeSubmitted: null as any,
    interviewed: false,
    codeAssessmentRequested: false,
  });
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();
  
  useEffect(() => {
    const saved = localStorage.getItem(`metrics-${id}`);
    if (saved) setMetrics(JSON.parse(saved));
  }, [id]);
  
  useEffect(() => {
    if (job && metrics.timeSubmitted) {
      const timeSub = new Date(metrics.timeSubmitted).getTime();
      const appliedOn = new Date(job.InitialApplicationDate).getTime();
      
      const timeSubDay = new Date(timeSub).setHours(0, 0, 0, 0);
      const appliedOnDay = new Date(appliedOn).setHours(0, 0, 0, 0);
      
      if (timeSubDay > appliedOnDay) {
        const correctedMetrics = { ...metrics, timeSubmitted: job.InitialApplicationDate as any };
        setMetrics(correctedMetrics);
        localStorage.setItem(`metrics-${id}`, JSON.stringify(correctedMetrics));
      }
    }
  }, [job, metrics.timeSubmitted, id]);
  
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
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="interviewed" checked={metrics.interviewed} onChange={handleMetricChange} />
              <span style={{ color: 'lightgreen' }}>Interviewed</span>
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
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="isEasyApply" checked={metrics.isEasyApply} onChange={handleMetricChange} />
              <span>Easy Apply</span>
            </label>
            <label className="checkbox-label" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
              <input type="checkbox" name="codeAssessmentRequested" checked={metrics.codeAssessmentRequested} onChange={handleMetricChange} />
              <span style={{ color: 'lightgreen' }}>Code Assessment Requested</span>
            </label>
            
            <div style={{ marginTop: '10px' }}>
              <label 
                htmlFor="timeSubmittedInput"
                style={{ display: 'block', marginBottom: '4px' }}
                >
                Time Submitted:
              </label>
              <input 
                type="datetime-local" 
                name="timeSubmitted" 
                value={
                  metrics.timeSubmitted 
                  ? new Date(new Date(metrics.timeSubmitted).getTime() - new Date(metrics.timeSubmitted).getTimezoneOffset() * 60000).toISOString().slice(0, 16)
                  : ""
                } 
                onChange={(e) => {
                  if (e.target.value) {
                    setMetrics(prev => ({ ...prev, timeSubmitted: new Date(e.target.value) as any }));
                  } else {
                    setMetrics(prev => ({ ...prev, timeSubmitted: null as any }));
                  }
                }}
                style={{ background: 'var(--bg-main)', color: 'var(--text-main)', border: '1px solid var(--border-color)', padding: '4px' }}
              />
            </div>
          </div>
        </div>
        <button onClick={saveMetrics} className="button primary-button" style={{ marginTop: '20px' }}>
        Save Metrics
        </button>
      </section>
    </div>
  );
};