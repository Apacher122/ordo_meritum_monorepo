import "@/assets/styles/Components/UI/TabHeaderControls.css";

import React from "react";

export const TABS = [
  "User Info",
  "Education",
  "Resume",
  "Cover Letter",
  "About Me",
] as const;

export type ProfileTab = typeof TABS[number];

interface TabHeaderControlsProps {
  activeTab: ProfileTab;
  onTabChange: (tab: ProfileTab) => void;
}

/**
 * Renders a set of tab buttons for navigating different sections of a form or page.
 * @param {TabHeaderControlsProps} props The props for the component.
 * @returns {React.FC<TabHeaderControlsProps>}
 */
export const TabHeaderControls: React.FC<TabHeaderControlsProps> = ({
  activeTab,
  onTabChange,
}) => {
  return (
    <div className="tab-bar">
      {TABS.map((tab) => (
        <button
          key={tab}
          className={`tab ${activeTab === tab ? "active" : ""}`}
          onClick={() => onTabChange(tab)}
        >
          {tab}
        </button>
      ))}
    </div>
  );
};