import '@/assets/styles/Components/UI/JobSelectorDropdown.css';

import { AppliedJob } from '../types';
import React from 'react';
import { useApplication } from '../providers/ApplicationProvider';

interface DualJobSelectorDropdownProps {
  jobsWithDoc: AppliedJob[];
  jobsWithoutDoc: AppliedJob[];
}

/**
 * A component that displays a dropdown selector
 * for the user to select a job. The dropdown selector is divided into two
 * groups: jobs with a document and jobs without a document.
 * 
 * @param {DualJobSelectorDropdownProps} props - The props object
 * @param {AppliedJob[]} props.jobsWithDoc - The list of jobs with a document
 * @param {AppliedJob[]} props.jobsWithoutDoc - The list of jobs without a document
 * @returns {React.ReactElement} - A React element representing the dropdown selector
 */
export const DualJobSelectorDropdown: React.FC<DualJobSelectorDropdownProps> = ({
  jobsWithDoc,
  jobsWithoutDoc,
}) => {
  const { selectedId, setSelectedId } = useApplication();

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
        {jobsWithDoc.length > 0 && (
          <optgroup label="Jobs with Document">
            {jobsWithDoc.map(job => (
              <option key={job.RoleID} value={job.RoleID}>
                {job.CompanyProperName} - {job.JobTitle}
              </option>
            ))}
          </optgroup>
        )}
        {jobsWithoutDoc.length > 0 && (
          <optgroup label="Jobs without Document">
            {jobsWithoutDoc.map(job => (
              <option key={job.RoleID} value={job.RoleID}>
                {job.CompanyProperName} - {job.JobTitle}
              </option>
            ))}
          </optgroup>
        )}
      </select>
    </div>
  );
};