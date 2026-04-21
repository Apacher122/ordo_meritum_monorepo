import { ApplicationStatus, AppliedJob } from '../types';
import React, { useEffect, useRef, useState } from 'react';

import { ApplicationListRow } from './ApplicationListRow';

interface ApplicationListViewProps {
    jobs: AppliedJob[];
    onStatusUpdate: (roleId: number, newStatus: ApplicationStatus) => void;
    onDateUpdate: (roleId: number, newDate: Date) => void;
    onDelete: (roleId: number) => void;
}

interface CachedRow {
  element: HTMLElement;
  top: number;
  height: number;
}

export const ApplicationListView: React.FC<ApplicationListViewProps> = ({ 
    jobs, 
    onStatusUpdate,
    onDateUpdate,
    onDelete 
}) => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const rowsCache = useRef<CachedRow[]>([]);
  const [visibleCount, setVisibleCount] = useState(50);

  const visibleJobs = jobs.slice(0, visibleCount);

  const applyScrollAnimations = () => {
    const container = scrollContainerRef.current;
    if (!container) return;

    const { scrollTop, clientHeight } = container;
    const rows = rowsCache.current

    if (rows.length === 0) return;

    for (const row of rows) {
      const { element, top, height } = row;
      const rowTop = top - scrollTop;

      let rotation = 0;
      let opacity = 1;
      let scale = 1;

      if (rowTop < 0) {
          const percentage = Math.min(1, Math.abs(rowTop) / height);
          rotation = -110 * percentage;
          opacity = 1 - percentage;
          scale = 1 - (0.3 * percentage);
      } else if (rowTop + height > clientHeight) {
          const distancePast = (rowTop + height) - clientHeight;
          const percentage = Math.min(1, distancePast / height);
          rotation = 110 * percentage;
          opacity = 1 - percentage;
          scale = 1 - (0.3 * percentage);
      }
      
      element.style.transform = `rotateX(${rotation}deg) scale(${scale})`;
      element.style.opacity = `${Math.max(0, opacity)}`;
    }
  };

  const buildAnimationCacheAndAnimate = () => {
    const container = scrollContainerRef.current;
    if (!container) return;

    const rowElements = Array.from(container.children).filter(
      (child) => child.classList.contains('application-row')
    ) as HTMLElement[];

    if (rowElements.length === 0) {
      rowsCache.current = [];
      return;
    }

    rowsCache.current = rowElements.map(element => ({
      element: element,
      top: element.offsetTop,
      height: element.offsetHeight,
    }));

    applyScrollAnimations();
  };
  
  useEffect(() => {
    const container = scrollContainerRef.current;
    if (!container) return;

    const handleScroll = () => applyScrollAnimations();
    
    container.addEventListener('scroll', handleScroll);
    
    const resizeObserver = new ResizeObserver(() => {
      buildAnimationCacheAndAnimate();
    });
    resizeObserver.observe(container);

    return () => {
      container.removeEventListener('scroll', handleScroll);
      resizeObserver.unobserve(container);
    };
  }, []);

  useEffect(() => {
    const timer = setTimeout(() => {
      buildAnimationCacheAndAnimate();
    }, 0);

    return () => clearTimeout(timer);
  }, [visibleJobs]);

  return (
    <div ref={scrollContainerRef} className="application-list-view">
      {visibleJobs.length > 0 ? (
        visibleJobs.map((app) => (
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
      
      {jobs.length > visibleCount && (
        <button 
           onClick={() => setVisibleCount(prev => prev + 50)}
           className="button"
           style={{ margin: '20px auto', display: 'block' }}
        >
          Load More
        </button>
      )}
    </div>
  );
};