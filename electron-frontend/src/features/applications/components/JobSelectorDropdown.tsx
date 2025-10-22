import '@/assets/styles/Components/UI/JobSelectorDropdown.css';

import { AppliedJob } from '../types';
import React from 'react';
import { useApplication } from '../providers/ApplicationProvider';

interface JobSelectorDropdownProps {
  jobs?: AppliedJob[];
}

export const JobSelectorDropdown: React.FC<JobSelectorDropdownProps> = ({ jobs: jobsProp }) => {
  const { jobs: jobsFromContext, selectedId, setSelectedId } = useApplication();

  const jobsToDisplay = jobsProp || jobsFromContext;

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
        {jobsToDisplay.map(job => (
          <option key={job.RoleID} value={job.RoleID}>
            {job.CompanyProperName} - {job.JobTitle}
          </option>
        ))}
      </select>
    </div>
  );
};