import { JobSelectorDropdown } from "./JobSelectorDropdown";
import React from "react";

export const DefaultHeaderControls: React.FC = () => {
  return (
    <div className="default-header-controls">
      <JobSelectorDropdown />
    </div>
  );
};
