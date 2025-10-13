import "../../assets/styles/components/Layouts/SideBar.css";

import React from "react";

/**
 * @interface SidebarItem
 * Defines the shape of a single item to be displayed in the sidebar.
 * @property {string} id - The unique identifier for the item, typically the URL path it navigates to.
 * @property {string} label - The text to display for the navigation item.
 */
interface SidebarItem {
  id: string;
  label: string;
}

/**
 * @interface SidebarProps
 * Defines the props for the Sidebar component.
 * @property {SidebarItem[]} items - An array of sidebar items to display.
 * @property {string} activeId - The ID of the currently active item, typically the current URL path.
 * @property {(id: string) => void} onSelect - The callback function to execute when an item is clicked.
 */
interface SidebarProps {
  items: SidebarItem[];
  activeId: string;
  onSelect: (id: string) => void;
}

/**
 * A reusable navigation sidebar component.
 * It is responsible for displaying navigation links and highlighting the active page.
 */
export const Sidebar: React.FC<SidebarProps> = ({
  items,
  activeId,
  onSelect,
}) => {
  return (
    <nav className="sidebar">
      <ul className="sidebar-list">
        {items.map(({ id, label }) => (
          <li
            key={id}
            className={`sidebar-item ${activeId === id ? "active" : ""}`}
            onClick={() => onSelect(id)}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                onSelect(id);
              }
            }}
            role="button"
            tabIndex={0}
          >
            {label}
          </li>
        ))}
      </ul>
    </nav>
  );
};

export default Sidebar;