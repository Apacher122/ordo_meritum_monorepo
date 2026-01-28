import '../../assets/styles/Components/UI/Header.css';

import React from 'react';
import { useHeaderContext } from './providers/HeaderProvider';

/**
 * Renders the main application header. It consumes the HeaderContext to display
 * a dynamic title, subtitle, and set of controls provided by the active page.
 * @returns {React.FC}
 */
export const Header: React.FC = () => {  
  const { title, subtitle, controls } = useHeaderContext();

  return (
    <header className="header-container">
      <div className="header-content">
        <div className="job-title-section">
          <h1 className="company-name">{title}</h1>
          <p className="job-position">{subtitle}</p>
        </div>
        <div className="header-controls">
          {controls}
        </div>
      </div>
    </header>
  );
};
