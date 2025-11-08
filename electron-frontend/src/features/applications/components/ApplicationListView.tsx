import { ApplicationStatus, AppliedJob } from '../types';
import React, { useEffect, useRef } from 'react';

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

/**
 * Renders a virtualized list of job applications with scroll-based animations.
 *
 * This component displays a list of `ApplicationListRow` items and applies a 3D
 * rotation, scale, and opacity effect to rows as they approach the top or bottom
 * edge of the viewport during scrolling.
 *
 * @performance
 * To achieve smooth scrolling performance and avoid layout thrashing, this
 * component implements a caching strategy:
 * 1.  It queries the DOM for row elements and caches their layout data
 * (`offsetTop`, `offsetHeight`) in a `useRef` (`rowsCache`).
 * 2.  This cache is built/rebuilt only when the `jobs` list changes or when the
 * container is resized.
 * 3.  The `onScroll` handler (`applyScrollAnimations`) runs on every scroll tick
 * but only performs fast calculations using the cached data. It *only*
 * writes to the DOM (updating `transform` and `opacity`), avoiding
 * expensive DOM reads in the scroll loop.
 *
 * @param {ApplicationListViewProps} props The component props.
 * @param {AppliedJob[]} props.jobs The array of job applications to display.
 * @param {(roleId: number, newStatus: ApplicationStatus) => void} props.onStatusUpdate
 * Callback function to update a job's status.
 * @param {(roleId: number, newDate: Date) => void} props.onDateUpdate
 * Callback function to update a job's application date.
 * @param {(roleId: number) => void} props.onDelete Callback function to delete a job.
 * @returns {React.ReactElement} The rendered ApplicationListView component.
 */
export const ApplicationListView: React.FC<ApplicationListViewProps> = ({ 
    jobs, 
    onStatusUpdate,
    onDateUpdate,
    onDelete 
}) => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const rowsCache = useRef<CachedRow[]>([]);

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