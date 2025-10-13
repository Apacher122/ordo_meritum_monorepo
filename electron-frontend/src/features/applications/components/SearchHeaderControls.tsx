import "@/assets/styles/Components/UI/SearchHeaderControls.css";

import React from "react";

interface SearchHeaderControlsProps {
  onSearch: (query: string) => void;
  initialQuery?: string;
}

export const SearchHeaderControls: React.FC<SearchHeaderControlsProps> = ({
  onSearch,
  initialQuery = "",
}) => {
  return (
    <div className="search-controls-container">
      <input
        type="search"
        placeholder="Search by company or title..."
        onChange={(e) => onSearch(e.target.value)}
        defaultValue={initialQuery}
        className="search-input"
        aria-label="Search job applications"
      />
    </div>
  );
};