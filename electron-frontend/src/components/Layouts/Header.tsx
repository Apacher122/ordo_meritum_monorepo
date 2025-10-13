import '../../assets/styles/Components/UI/Header.css';

import React from 'react';
import { useHeaderContext } from './providers/HeaderProvider';

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
