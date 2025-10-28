import '@/assets/styles/Components/Layouts/ApplicationListRow.css';

import { ApplicationStatus, AppliedJob } from '../types';
import React, { useState } from 'react';

interface ApplicationListRowProps {
  application: AppliedJob;
  onStatusUpdate: (roleId: number, newStatus: ApplicationStatus) => void;
  onDateUpdate: (roleId: number, newDate: Date) => void;
  onDelete: (roleId: number) => void;
  className?: string;
}

const statusOptions: ApplicationStatus[] = ['Rejected', 'Offered', 'Open', 'Closed', 'Moved', 'Not applied', 'Ghosted', 'Interviewing'];

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
  const [isUpdating, setIsUpdating] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isEditingDate, setIsEditingDate] = useState(false);

  const handleStatusChange = async (event: React.ChangeEvent<HTMLSelectElement>) => {
    const newFrontendStatus = event.target.value as ApplicationStatus;
    setIsUpdating(true);
    await onStatusUpdate(application.RoleID, newFrontendStatus);
    setIsUpdating(false);
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

        {showDeleteConfirm && (
            <ConfirmationModal
                message="This will permanently remove this application from your list."
                onConfirm={handleDelete}
                onCancel={() => setShowDeleteConfirm(false)}
            />
        )}
    </div>
  );
};