import "@/assets/styles/components/UI/TabHeaderControls.css";

import React from "react";

export type ProfileTab =
  | "User Info"
  | "Education"
  | "Resume"
  | "Cover Letter"
  | "About Me";

const TABS: ProfileTab[] = [
  "User Info",
  "Education",
  "Resume",
  "Cover Letter",
  "About Me",
];

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