import { JobSelectorDropdown } from "./JobSelectorDropdown";
import React from "react";

/**
 * A component that composes the default controls for the main header,
 * primarily featuring the job selection dropdown.
 */
export const DefaultHeaderControls: React.FC = () => {
  return (
    <div className="default-header-controls">
      <JobSelectorDropdown />
    </div>
  );
};
