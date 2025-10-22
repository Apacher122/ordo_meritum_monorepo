import { ApplicationStatus, AppliedJob } from '../types';
import React, { useEffect, useRef } from 'react';

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
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = scrollContainerRef.current;
    if (!container) return;

    const handleScroll = () => {
      const { scrollTop, clientHeight } = container;
      const rows = Array.from(container.children).filter(
        (child) => child.classList.contains('application-row')
      ) as HTMLElement[];

      if (rows.length === 0) return;

      const rowHeight = rows[0].offsetHeight;
      if (rowHeight === 0) return;

      for (const row of rows) {
        const rowTop = row.offsetTop - scrollTop;
        
        let rotation = 0;
        let opacity = 1;
        let scale = 1;

        if (rowTop < 0) {
            const percentage = Math.min(1, Math.abs(rowTop) / rowHeight);
            rotation = -110 * percentage;
            opacity = 1 - percentage;
            scale = 1 - (0.3 * percentage);
        } else if (rowTop + rowHeight > clientHeight) {
            const distancePast = (rowTop + rowHeight) - clientHeight;
            const percentage = Math.min(1, distancePast / rowHeight);
            rotation = 110 * percentage;
            opacity = 1 - percentage;
            scale = 1 - (0.3 * percentage);
        }
        
        row.style.transform = `rotateX(${rotation}deg) scale(${scale})`;
        row.style.opacity = `${Math.max(0, opacity)}`;
      }
    };
    
    handleScroll();
    container.addEventListener('scroll', handleScroll);
    
    const resizeObserver = new ResizeObserver(handleScroll);
    resizeObserver.observe(container);

    return () => {
      container.removeEventListener('scroll', handleScroll);
      resizeObserver.unobserve(container);
    };
  }, [jobs]);

  return (
    <div ref={scrollContainerRef} className="application-list-view">
      {jobs.length > 0 ? (
        jobs.map((app) => (
          <ApplicationListRow
            key={app.RoleID}
            application={app}
            onStatusUpdate={onStatusUpdate}
            onDateUpdate={onDateUpdate}
            onDelete={onDelete}
            className="application-row"
          />
        ))
      ) : (
        <p>No matching applications found.</p>
      )}
    </div>
  );
};