import '@/assets/styles/Components/Layouts/ApplicationMetrics.css';

import React, { useMemo } from 'react';

import { AppliedJob } from '../types';

export interface ApplicationMetricsData {
  applicationsSent: number;
  rejections: number;
  ghosts: number;
  stillOpen: number;
  interviewing: number;
  offers: number;
  interviewRate: number;
  rejectionRate: number;
  ghostedRate: number;
  openAppsRate: number;
  appsSentToday: number;
  avgAppsPerDay: string;
}

interface ApplicationMetricsProps {
  jobs: AppliedJob[];
  onClose: () => void; 
}

/**
 * Renders a modal overlay displaying key statistics derived from the user's job applications.
 *
 * This component is designed to be displayed on-demand (e.g., when a user clicks
 * a "Show Metrics" button). It receives the complete, unfiltered list of job
 * applications as a prop.
 *
 * All statistical calculations (e.g., rejection rates, totals, averages)
 * are performed internally within this component. This logic is memoized
 * using the `useMemo` hook, ensuring that the expensive calculations
 * only re-run if the `jobs` prop array changes. This optimizes performance
 * by deferring computation until it is actually needed.
 *
 * @param {ApplicationMetricsProps} props The component props.
 * @param {AppliedJob[]} props.jobs The complete array of job applications to be analyzed.
 * @param {() => void} props.onClose A callback function to be invoked when the modal's
 * close button is clicked.
 * @returns {React.ReactElement} The rendered ApplicationMetrics modal component.
 */
export const ApplicationMetrics: React.FC<ApplicationMetricsProps> = ({ jobs, onClose }) => {
  const metrics = useMemo((): ApplicationMetricsData => {
    const appsSent = jobs.filter((j) => j.ApplicationStatus !== "Not applied");
    const applicationsSent = appsSent.length;
    if (applicationsSent === 0) {
      return {
        applicationsSent: 0,
        rejections: 0,
        ghosts: 0,
        stillOpen: 0,
        interviewing: 0,
        offers: 0,
        interviewRate: 0,
        rejectionRate: 0,
        ghostedRate: 0,
        openAppsRate: 0,
        appsSentToday: 0,
        avgAppsPerDay: "0.0",
      };
    }

    const rejections = appsSent.filter(
      (j) => j.ApplicationStatus === "Rejected"
    ).length;
    const ghosts = appsSent.filter(
      (j) => j.ApplicationStatus === "Ghosted"
    ).length;
    const stillOpen = appsSent.filter(
      (j) =>
        j.ApplicationStatus === "Open" || j.ApplicationStatus === "Interviewing"
    ).length;
    const interviewing = appsSent.filter(
      (j) => j.ApplicationStatus === "Interviewing"
    ).length;
    const offers = appsSent.filter(
      (j) => j.ApplicationStatus === "Offered"
    ).length;

    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const appsSentToday = appsSent.filter((j) => {
      const appDate = new Date(j.InitialApplicationDate);
      appDate.setHours(0, 0, 0, 0);
      return appDate.getTime() === today.getTime();
    }).length;

    const firstAppDate = appsSent.reduce((oldest, job) => {
      const jobDate = new Date(job.InitialApplicationDate);
      return jobDate.getTime() < oldest.getTime() ? jobDate : oldest;
    }, new Date());
    const daysSinceFirstApp = Math.max(
      1,
      Math.ceil((Date.now() - firstAppDate.getTime()) / (1000 * 60 * 60 * 24))
    );
    const avgAppsPerDay = (applicationsSent / daysSinceFirstApp).toFixed(1);

    return {
      applicationsSent,
      rejections,
      ghosts,
      stillOpen,
      interviewing,
      offers,
      interviewRate: (interviewing / applicationsSent) * 100,
      rejectionRate: (rejections / applicationsSent) * 100,
      ghostedRate: (ghosts / applicationsSent) * 100,
      openAppsRate: (stillOpen / applicationsSent) * 100,
      appsSentToday,
      avgAppsPerDay,
    };
  }, [jobs]);
  return (
    <div className="metrics-modal-overlay">
      <div className="metrics-modal">
        <div className="metrics-header">
          <h2>Application Metrics</h2>
          <button onClick={onClose} className="close-button">&times;</button>
        </div>
        <div className="metrics-grid">
          <div className="metric-item">
            <span className="metric-value">{metrics.applicationsSent}</span>
            <span className="metric-label">Total Sent</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.rejections}</span>
            <span className="metric-label">Rejections</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.ghosts}</span>
            <span className="metric-label">Ghosted</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.stillOpen}</span>
            <span className="metric-label">Still Open</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.interviewing}</span>
            <span className="metric-label">Interviewing</span>
          </div>
           <div className="metric-item">
            <span className="metric-value">{metrics.offers}</span>
            <span className="metric-label">Offers</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.interviewRate.toFixed(1)}%</span>
            <span className="metric-label">Interview Rate</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.rejectionRate.toFixed(1)}%</span>
            <span className="metric-label">Rejection Rate</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.ghostedRate.toFixed(1)}%</span>
            <span className="metric-label">Ghosted Rate</span>
          </div>
           <div className="metric-item">
            <span className="metric-value">{metrics.openAppsRate.toFixed(1)}%</span>
            <span className="metric-label">Open App Rate</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.appsSentToday}</span>
            <span className="metric-label">Sent Today</span>
          </div>
          <div className="metric-item">
            <span className="metric-value">{metrics.avgAppsPerDay}</span>
            <span className="metric-label">Avg Apps / Day</span>
          </div>
        </div>
      </div>
    </div>
  );
};
