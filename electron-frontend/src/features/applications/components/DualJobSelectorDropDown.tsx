import '@/assets/styles/Components/UI/JobSelectorDropdown.css';

import { AppliedJob } from '../types';
import React from 'react';
import { useApplication } from '../providers/ApplicationProvider';

interface DualJobSelectorDropdownProps {
  jobsWithDoc: AppliedJob[];
  jobsWithoutDoc: AppliedJob[];
  id?: string;
}


/**
 * A dropdown selector component that displays two groups of job applications.
 * The first group displays job applications that have a document associated with them.
 * The second group displays job applications that do not have a document associated with them.
 *
 * @param {DualJobSelectorDropdownProps} props - The component props.
 * @param {AppliedJob[]} props.jobsWithDoc - The array of job applications with documents.
 * @param {AppliedJob[]} props.jobsWithoutDoc - The array of job applications without documents.
 * @param {string} [props.id] - The id of the dropdown selector. Defaults to an empty string.
 * @returns {React.ReactElement} - The rendered DualJobSelectorDropdown component.
 */
export const DualJobSelectorDropdown: React.FC<DualJobSelectorDropdownProps> = ({
  jobsWithDoc,
  jobsWithoutDoc,
  id,
}) => {
  const { selectedId, setSelectedId } = useApplication();

  const handleSelectionChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const newId = Number.parseInt(event.target.value, 10);
    setSelectedId(Number.isNaN(newId) ? null : newId);
  };

  return (
    <div className="job-selector-container">
      <select
        id={id}
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