import '../Styles/TextBox.css';

import React from 'react';

interface LargeTextBoxProps {
  value: string;
  onChange: (newValue: string) => void;
  onSave: () => void;
  placeholder?: string;
  disabled?: boolean;
}

export default function LargeTextBox({ value, onChange, onSave, placeholder, disabled }: LargeTextBoxProps) {
  return (
    <div className="large-textbox-container">
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder || 'Paste LinkedIn job description here...'}
        className="large-textbox-textarea"
        disabled={disabled}
      />
      <button
        onClick={onSave}
        className="large-textbox-button"
        disabled={disabled || !value.trim()}
      >
        Save
      </button>
    </div>
  );
}
