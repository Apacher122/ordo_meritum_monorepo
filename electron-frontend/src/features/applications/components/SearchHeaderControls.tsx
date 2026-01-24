import "@/assets/styles/Components/UI/SearchHeaderControls.css";

import React from "react";

interface SearchHeaderControlsProps {
  onSearch: (query: string) => void;
  value: string;
}

/**
 * A component for searching job applications by company or title.
 * @param {SearchHeaderControlsProps} props - The props for the SearchHeaderControls component.
 * @param {function} props.onSearch - A callback function that is called when the user submits the search query.
 * @param {string} props.value - The current value of the search input field.
 * @returns {React.ReactElement} - A React element representing the SearchHeaderControls component.
 */
export const SearchHeaderControls: React.FC<SearchHeaderControlsProps> = ({
  onSearch,
  value,
}) => {
  return (
    <div className="search-controls-container">
      <input
        type="search"
        placeholder="Search by company or title..."
        onChange={(e) => onSearch(e.target.value)}
        value={value}
        className="search-input"
        aria-label="Search job applications"
      />
    </div>
  );
};