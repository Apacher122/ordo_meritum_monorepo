import '@/assets/styles/Components/UI/JobSelectorDropdown.css';

import { AppliedJob } from '../types';
import React from 'react';
import { useApplication } from '../providers/ApplicationProvider';

interface DualJobSelectorDropdownProps {
  jobsWithDoc: AppliedJob[];
  jobsWithoutDoc: AppliedJob[];
}

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