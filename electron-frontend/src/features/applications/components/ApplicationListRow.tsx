import '@/assets/styles/Components/Layouts/ApplicationListRow.css';

import { ApplicationStatus, AppliedJob, statusOptions } from '../types';
import React, { useEffect, useState } from 'react';

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

/**
 * A modal component that asks the user to confirm an action.
 *
 * @param {string} message The message to be displayed to the user.
 * @param {() => void} onConfirm The function to be called when the user confirms the action.
 * @param {() => void} onCancel The function to be called when the user cancels the action.
 * @returns {React.FC} A React component that displays a confirmation modal.
 */
const ConfirmationModal: React.FC<{ message: string; onConfirm: () => void; onCancel: () => void }> = ({ message, onConfirm, onCancel }) => (
    <div className="confirmation-modal-overlay">
        <div className="confirmation-modal">
            <h2>Are you sure?</h2>
            <p>{message}</p>
            <div className="confirmation-modal-buttons">
                <button onClick={onCancel} className="cancel-button">Cancel</button>
                <button onClick={onConfirm} className="confirm-button">Confirm</button>
            </div>
        </div>
    </div>
);

/**
 * A component that displays a row of information about an application job.
 *
 * @param {AppliedJob} application The application job to be displayed.
 * @param {(roleId: number, newStatus: ApplicationStatus) => void} onStatusUpdate A function to be called when the user updates the status of the application.
 * @param {(roleId: number, newDate: Date) => void} onDateUpdate A function to be called when the user updates the date of the application.
 * @param {(roleId: number) => void} onDelete A function to be called when the user deletes the application.
 * @param {string} className An optional class name to be applied to the component.
 */
export const ApplicationListRow: React.FC<ApplicationListRowProps> = ({ application, onStatusUpdate, onDateUpdate, onDelete, className }) => {
  const navigate = useNavigate()
	const { setSelectedId } = useApplication();
  const [isUpdating, setIsUpdating] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isEditingDate, setIsEditingDate] = useState(false);
	const [hasResume, setHasResume] = useState(false);
  const [hasCoverLetter, setHasCoverLetter] = useState(false);

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
  
  const handleDelete = () => {
      onDelete(application.RoleID);
      setShowDeleteConfirm(false);
  };

  const handleDateChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const newDate = new Date(event.target.value);
    const timezoneOffset = newDate.getTimezoneOffset() * 60000;
    const adjustedDate = new Date(newDate.getTime() + timezoneOffset);
    setIsEditingDate(false);
    await onDateUpdate(application.RoleID, adjustedDate);
  };

  const formatDate = (date: Date) => {
    if (!date || Number.isNaN(date.getTime())) return 'N/A';
    return date.toLocaleDateString();
  };
  
  const formatDateForInput = (date: Date): string => {
      if (!date || Number.isNaN(date.getTime())) return '';
      return date.toISOString().split('T')[0];
  };

  const rowClasses = `${className || ''} application-row ${application.ApplicationStatus === 'Rejected' ? 'rejected' : ''}`.trim();

  return (
    <div className={rowClasses}>
        <div className="application-content">
            <div className="application-info">
                <div className="company-position">
                    <span>{application.CompanyProperName},</span>
                    <span className="position">{application.JobTitle}</span>
                </div>
                 <div className="status-selector">
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
                            onClick={() => handleOpenDocument('resume')}
                            title="View Resume"
                        >
                            Resume
                        </button>
                    )}
                    {hasCoverLetter && (
                        <button 
                            className="doc-shortcut-button" 
                            onClick={() => handleOpenDocument('cover-letter')}
                            title="View Cover Letter"
                        >
                            Cover Letter
                        </button>
                    )}
                </div>
            </div>

            <div className="application-actions">
                <div className="applied-date-container">
                    {isEditingDate ? (
                        <input 
                            type="date"
                            defaultValue={formatDateForInput(application.InitialApplicationDate)}
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
                            <i>Applied On: {formatDate(application.InitialApplicationDate)}</i>
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
};