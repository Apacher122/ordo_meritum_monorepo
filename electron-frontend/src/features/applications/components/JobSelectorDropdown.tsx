import '@/assets/styles/Components/UI/JobSelectorDropdown.css';

import React from 'react';
import { useApplication } from '../providers/ApplicationProvider';

/**
 * A dropdown component that allows the user to select a job application.
 * It reads the list of jobs and the currently selected ID from the ApplicationProvider,
 * and updates the selected ID when the user makes a new selection.
 */
export const JobSelectorDropdown: React.FC = () => {
  const { jobs, selectedId, setSelectedId } = useApplication();

  const handleSelectionChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const newId = Number.parseInt(event.target.value, 10);
    setSelectedId(Number.isNaN(newId) ? null : newId);
  };

  return (
    <div className="job-selector-container">
      <select 
        className="job-selector-dropdown"
        value={selectedId ?? ''} 
        onChange={handleSelectionChange}
      >
        <option value="">Select a Job...</option>
        {jobs.map(job => (
          <option key={job.RoleID} value={job.RoleID}>
            {job.CompanyProperName} - {job.JobTitle}
          </option>
        ))}
      </select>
    </div>
  );
};
