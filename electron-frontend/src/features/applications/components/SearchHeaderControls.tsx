import "@/assets/styles/Components/UI/SearchHeaderControls.css";

import React, { useEffect, useState } from "react";

import { useDebounce } from "../hooks/useDebounce";

interface SearchHeaderControlsProps {
  onSearch: (query: string) => void;
}

export const SearchHeaderControls: React.FC<SearchHeaderControlsProps> = ({ onSearch }) => {
  const [localValue, setLocalValue] = useState("");
  const debouncedSearchTerm = useDebounce(localValue, 300);

  useEffect(() => {
    onSearch(debouncedSearchTerm);
  }, [debouncedSearchTerm, onSearch]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setLocalValue(e.target.value);
  };

  return (
    <div className="search-controls-container">
      <input
        type="search"
        placeholder="Search by company or title..."
        onChange={handleChange}
        value={localValue}
        className="search-input"
        aria-label="Search job applications"
      />
    </div>
  );
};