import '@/assets/styles/Components/Layouts/ApplicationMetrics.css';

import React from 'react';

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
  metrics: ApplicationMetricsData;
  onClose: () => void; 
}

export const ApplicationMetrics: React.FC<ApplicationMetricsProps> = ({ metrics, onClose }) => {
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
