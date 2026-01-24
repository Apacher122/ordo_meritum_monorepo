import "@/assets/styles/Views/MatchSummaryView.css";

import React, { useEffect, useMemo, useState } from "react";
import {
  useSetHeaderControls,
  useSetHeaderSubtitle,
  useSetHeaderTitle,
} from "@/components/Layouts/providers/HeaderProvider";

import { AppliedJob } from "@/features/applications/types";
import { CircleProgress } from "@/components/UI/loaders";
import { JobSelectorDropdown } from "@/features/applications/components/JobSelectorDropdown";
import RadarChartOverview from "@/components/UI/charts/RadarChart";
import { SummaryInfo } from "@/shared/types";
import { useApplication } from "@/features/applications/providers/ApplicationProvider";
import { useDocumentManager } from "@/features/documents/hooks/useDocumentManager";
import { useMatchSummary } from "../hooks/useMatchSummary";

const getSummaryItemClassName = (summary: SummaryInfo) => {
  if (summary.summary_temperature > 0) return "summary-item good";
  if (summary.summary_temperature < 0) return "summary-item bad";
  return "summary-item neutral";
};

const getSummarySymbol = (temperature: number) => {
  if (temperature > 0) return '✓';
  if (temperature < 0) return '✗';
  return '-';
};

/**
 * Renders the Job Match Summary page. This page displays a detailed analysis
 * of how a candidate's resume matches a selected job posting, including an
 * overall score, detailed metrics, and AI-generated summaries.
 * @returns {React.FC} The MatchSummaryPage component.
 */
export const MatchSummaryPage: React.FC = () => {
  const { jobs, selectedJob } = useApplication();
  const { matchSummary, loading, error, getMatchSummary, hasLocalSummary } = useMatchSummary();
  const [jobsWithResume, setJobsWithResume] = useState<AppliedJob[]>([]);
  const { doesFileExist } = useDocumentManager(null, "", "", "resume");

  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  useEffect(() => {
    const findJobsWithResumes = async () => {
      const checkedJobs = await Promise.all(
        jobs.map(async (job) => {
          const exists = await doesFileExist(job.RoleID, "resume", job.CompanyName, job.JobTitle);
          return exists ? job : null;
        })
      );
      setJobsWithResume(checkedJobs.filter((job): job is AppliedJob => job !== null));
    };

    if (jobs.length > 0) findJobsWithResumes();
  }, [jobs, doesFileExist]);

  const headerControls = useMemo(() => (
    <div className="match-summary-header-controls">
      <JobSelectorDropdown jobs={jobsWithResume} />
      {selectedJob && !hasLocalSummary && (
        <button onClick={getMatchSummary} className="button" disabled={loading}>
          {loading ? "Generating..." : "Get New Match Summary"}
        </button>
      )}
    </div>
  ), [jobsWithResume, selectedJob, hasLocalSummary, getMatchSummary, loading]);

  useEffect(() => {
    if (selectedJob) {
      setHeaderTitle(selectedJob.CompanyProperName);
      setHeaderSubtitle(selectedJob.JobTitle);
    } else {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
    }
    setHeaderControls(headerControls);
  }, [selectedJob, headerControls, setHeaderTitle, setHeaderSubtitle, setHeaderControls]);

  if (loading) {
    return (
      <div className="page-content-placeholder">
        <CircleProgress />
        <p>Loading Match Summary...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="page-content-placeholder error-message">
        <h2>Error</h2>
        <p>{error}</p>
      </div>
    );
  }

  if (!selectedJob) {
    return (
      <div className="page-content-placeholder">
        Please select a job with a generated resume to view the match summary.
      </div>
    );
  }

  if (!matchSummary) {
    return (
      <div className="page-content-placeholder">
        <h2>No Match Summary Available</h2>
        <p>You can generate one using the button in the header.</p>
      </div>
    );
  }
  
  const {
    should_apply,
    should_apply_reasoning,
    overall_match_summary,
    metrics,
  } = matchSummary.match_summary;
  
  const shouldApplyClass = should_apply ? "should-apply-yes" : "should-apply-no";

  return (
    <div className="match-summary-view">
      <div className="main-summary-card">
        <div className="circular-container">
          <CircleProgress
            percentage={overall_match_summary.overall_match_score}
            size={150}
            isGenerating={false}
          />
        </div>
        <div className="summary-text">
          <h2 className="overall-match-score">Overall Match: {overall_match_summary.overall_match_score}%</h2>
          <h3 className={shouldApplyClass}>Should Apply? {should_apply ? "Yes" : "No"}</h3>
          <p>{should_apply_reasoning}</p>
          <ul className="summary-list">
            {overall_match_summary.summary.map((summary, index) => (
              <li key={`${summary.summary_text}-${index}`} className={getSummaryItemClassName(summary)}>
                <span className="symbol">{getSummarySymbol(summary.summary_temperature)}</span>
                <span className="text">{summary.summary_text}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
      
      <div className="radar-and-details">
        <div className="radar-chart-container">
          <RadarChartOverview metrics={metrics} />
        </div>
        <div className="match-detail-container">
          <div className="match-detail-scroll">
            {metrics.map((metric) => (
              <div key={metric.score_title} className="match-detail-card">
                <h3>{metric.score_title}</h3>
                <p>
                  <strong>Score:</strong> {metric.weighted_score}%
                </p>
                <p>
                  <strong>Reasoning:</strong> {metric.score_reason}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};