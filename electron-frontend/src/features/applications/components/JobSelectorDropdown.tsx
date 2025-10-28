import '@/assets/styles/Components/UI/JobSelectorDropdown.css';

import { AppliedJob } from '../types';
import React from 'react';
import { useApplication } from '../providers/ApplicationProvider';

interface JobSelectorDropdownProps {
  jobs?: AppliedJob[];
}

/**
 * A component that displays a dropdown selector
 * for the user to select a job.
 * 
 * @param {JobSelectorDropdownProps} props - The props object
 * @param {AppliedJob[]} [props.jobs] - The list of jobs to display in the dropdown selector
 * @returns {React.ReactElement} - A React element representing the dropdown selector
 */

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