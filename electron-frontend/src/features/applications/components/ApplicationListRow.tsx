import { ApplicationStatus, AppliedJob, statusOptions } from '../types';
import React, { memo, useEffect, useMemo, useState } from 'react';

import { ConfirmationModal } from '@/components/UI';
import { createPortal } from 'react-dom';
import { useApplication } from '@/app/appProviders';
import { useDocumentManager } from '@/features/documents/hooks';
import { useNavigate } from 'react-router-dom';

interface ApplicationListRowProps {
  application: AppliedJob;
  onStatusUpdate: (roleId: number, newStatus: ApplicationStatus) => void;
  onDateUpdate: (roleId: number, newDate: Date) => void;
  onDelete: (roleId: number) => void;
  className?: string;
}

export const ApplicationListRow: React.FC<ApplicationListRowProps> = memo(({ application, onStatusUpdate, onDateUpdate, onDelete, className }) => {
  const navigate = useNavigate()
  const { setSelectedId } = useApplication();
  const [isUpdating, setIsUpdating] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isEditingDate, setIsEditingDate] = useState(false);
  const [hasResume, setHasResume] = useState(false);
  const [hasCoverLetter, setHasCoverLetter] = useState(false);
  const [metrics, setMetrics] = useState<any>({});

  useEffect(() => {
    const saved = localStorage.getItem(`metrics-${application.RoleID}`);
    if (saved) {
      setMetrics(JSON.parse(saved));
    }
  }, [application.RoleID]);

  useEffect(() => {
    if (metrics.timeSubmitted && application.InitialApplicationDate) {
      const timeSub = new Date(metrics.timeSubmitted).getTime();
      const appliedOn = new Date(application.InitialApplicationDate).getTime();
      const timeSubDay = new Date(timeSub).setHours(0, 0, 0, 0);
      const appliedOnDay = new Date(appliedOn).setHours(0, 0, 0, 0);
      
      if (timeSubDay > appliedOnDay) {
        const updatedMetrics = { ...metrics, timeSubmitted: application.InitialApplicationDate };
        setMetrics(updatedMetrics);
        localStorage.setItem(`metrics-${application.RoleID}`, JSON.stringify(updatedMetrics));
      }
    }
  }, [metrics.timeSubmitted, application.InitialApplicationDate, application.RoleID]);

  const displayDate = useMemo(() => {
    if (metrics.timeSubmitted) {
      const parsedDate = new Date(metrics.timeSubmitted);
      if (!Number.isNaN(parsedDate.getTime())) return parsedDate;
    }
    return application.InitialApplicationDate;
  }, [metrics.timeSubmitted, application.InitialApplicationDate]);

  const { doesFileExist } = useDocumentManager(
    application.RoleID,
    application.CompanyName,
    application.JobTitle,
    "resume"
  );

  useEffect(() => {
    const checkDocs = async () => {
      const resumeExists = await doesFileExist(application.RoleID, "resume", application.CompanyName, application.JobTitle);
      const clExists = await doesFileExist(application.RoleID, "cover-letter", application.CompanyName, application.JobTitle);
      setHasResume(resumeExists);
      setHasCoverLetter(clExists);
    };
    checkDocs();
  }, [application, doesFileExist]);
  
  const handleStatusChange = async (event: React.ChangeEvent<HTMLSelectElement>) => {
    const newFrontendStatus = event.target.value as ApplicationStatus;
    setIsUpdating(true);
    await onStatusUpdate(application.RoleID, newFrontendStatus);
    setIsUpdating(false);
  };

  const handleOpenDocument = (docType: 'resume' | 'cover-letter') => {
    setSelectedId(application.RoleID);
    navigate('/documents', { state: { initialDocType: docType } });
  };

  const handleRowClick = () => {
    setSelectedId(application.RoleID);
    navigate(`/applications/${application.RoleID}`);
  };
  
  const handleDelete = () => {
      onDelete(application.RoleID);
      setShowDeleteConfirm(false);
  };

  const handleDateChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.value) {
      const newDate = new Date(event.target.value);
      setIsEditingDate(false);
      await onDateUpdate(application.RoleID, newDate);
      
      const updatedMetrics = { ...metrics, timeSubmitted: newDate };
      setMetrics(updatedMetrics);
      localStorage.setItem(`metrics-${application.RoleID}`, JSON.stringify(updatedMetrics));
    } else {
      setIsEditingDate(false);
      
      const updatedMetrics = { ...metrics, timeSubmitted: null };
      setMetrics(updatedMetrics);
      localStorage.setItem(`metrics-${application.RoleID}`, JSON.stringify(updatedMetrics));
    }
  };

  const handleMetricChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, checked } = event.target;
    const updatedMetrics = { ...metrics, [name]: checked };
    setMetrics(updatedMetrics);
    localStorage.setItem(`metrics-${application.RoleID}`, JSON.stringify(updatedMetrics));
  };

  const formatDate = (date: Date) => {
    if (!date || Number.isNaN(date.getTime())) return 'N/A';
    return date.toLocaleString([], { year: 'numeric', month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' });
  };
  
  const formatDateForInput = (date: Date): string => {
      if (!date || Number.isNaN(date.getTime())) return '';
      const offsetMs = date.getTimezoneOffset() * 60000;
      const localDate = new Date(date.getTime() - offsetMs);
      return localDate.toISOString().slice(0, 16);
  };

  const rowClasses = `${className || ''} application-row ${application.ApplicationStatus === 'Rejected' ? 'rejected' : ''}`.trim();

  return (
    <div className={rowClasses}>
      <div className="application-content" onClick={handleRowClick} style={{ cursor: 'pointer' }}>
        <div className="application-info">
          <div className="company-position">
            <span>{application.CompanyProperName},</span>
            <span className="position">{application.JobTitle}</span>
          </div>
          <div className="status-selector" onClick={(e) => e.stopPropagation()}>
            <label htmlFor={`status-${application.RoleID}`}>Status:</label>
            <select
              id={`status-${application.RoleID}`}
                value={application.ApplicationStatus}
                onChange={handleStatusChange}
                disabled={isUpdating}
                className="status-dropdown"
            >
              {statusOptions.map(status => (
                  <option key={status} value={status}>{status}</option>
              ))}
            </select>
            {hasResume && (
              <button 
                className="doc-shortcut-button" 
                onClick={(e) => { e.stopPropagation(); handleOpenDocument('resume'); }}
                title="View Resume"
              >
                Resume
              </button>
            )}
            {hasCoverLetter && (
              <button 
                className="doc-shortcut-button" 
                onClick={(e) => { e.stopPropagation(); handleOpenDocument('cover-letter'); }}
                title="View Cover Letter"
              >
                Cover Letter
              </button>
            )}
          </div>
        </div>

        <div className="application-actions" onClick={(e) => e.stopPropagation()}>
          <div className="applied-date-container">
            {isEditingDate ? (
                <input 
                  type="datetime-local"
                  defaultValue={formatDateForInput(displayDate)}
                  onBlur={handleDateChange} 
                  autoFocus
                  className="date-input"
                />
            ) : (
              <button
                className="applied-date editable text-button"
                onClick={() => setIsEditingDate(true)}
                title="Click to edit date"
              >
                <i>Applied On: {formatDate(displayDate)}</i>
              </button>
            )}
          </div>
          <button 
            onClick={() => setShowDeleteConfirm(true)}
            className="remove-button"
            title="Remove Application"
          >
            🗑️
          </button>
        </div>
      </div>

      <div 
        className="quick-metrics-bar" 
        onClick={(e) => e.stopPropagation()} 
        style={{ 
          padding: '8px 15px', 
          display: 'flex', 
          gap: '15px', 
          flexWrap: 'wrap', 
          fontSize: '0.85em', 
          borderTop: '1px solid var(--border-color)', 
          background: 'var(--bg-secondary)',
          borderBottomLeftRadius: 'inherit',
          borderBottomRightRadius: 'inherit'
        }}
      >
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
          <input type="checkbox" name="hasSummary" checked={!!metrics.hasSummary} onChange={handleMetricChange} />
          <span>Summary</span>
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
          <input type="checkbox" name="hasCoverLetter" checked={!!metrics.hasCoverLetter} onChange={handleMetricChange} />
          <span>Cover Letter</span>
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
          <input type="checkbox" name="isVeteran" checked={!!metrics.isVeteran} onChange={handleMetricChange} />
          <span>Veteran</span>
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
          <input type="checkbox" name="isDisabled" checked={!!metrics.isDisabled} onChange={handleMetricChange} />
          <span>Disability</span>
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
          <input type="checkbox" name="raceMentioned" checked={!!metrics.raceMentioned} onChange={handleMetricChange} />
          <span>Race Mentioned</span>
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
          <input type="checkbox" name="isEasyApply" checked={!!metrics.isEasyApply} onChange={handleMetricChange} />
          <span>Easy Apply</span>
        </label>
        <span style={{ borderLeft: '1px solid var(--border-color)', margin: '0 5px' }}></span>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer', color: 'lightgreen' }}>
          <input type="checkbox" name="interviewed" checked={!!metrics.interviewed} onChange={handleMetricChange} />
          <span>Interviewed</span>
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer', color: 'lightgreen' }}>
          <input type="checkbox" name="codeAssessmentRequested" checked={!!metrics.codeAssessmentRequested} onChange={handleMetricChange} />
          <span>Code Assessment</span>
        </label>
      </div>

      {showDeleteConfirm && createPortal(
        <ConfirmationModal
          message="This will permanently remove this application from your list."
          onConfirm={handleDelete}
          onCancel={() => setShowDeleteConfirm(false)}
        />,
        document.body
      )}
    </div>
  );
});

ApplicationListRow.displayName = 'ApplicationListRow';