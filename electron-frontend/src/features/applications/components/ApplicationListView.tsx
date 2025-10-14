import { ApplicationStatus, AppliedJob } from '../types';

import { ApplicationListRow } from './ApplicationListRow';

interface ApplicationListViewProps {
    jobs: AppliedJob[];
    onStatusUpdate: (roleId: number, newStatus: ApplicationStatus) => void;
    onDateUpdate: (roleId: number, newDate: Date) => void;
    onDelete: (roleId: number) => void;
}

export const ApplicationListView: React.FC<ApplicationListViewProps> = ({ 
    jobs, 
    onStatusUpdate,
    onDateUpdate,
    onDelete 
}) => {
  return (
    <div className="application-list-view">
      {jobs.length > 0 ? (
        jobs.map((app) => (
          
          <ApplicationListRow
            key={app.RoleID}
            application={app}
            onStatusUpdate={onStatusUpdate}
            onDateUpdate={onDateUpdate}
            onDelete={onDelete}
          />
        ))
      ) : (
        <p>No matching applications found.</p>
      )}
    </div>
  );
};

