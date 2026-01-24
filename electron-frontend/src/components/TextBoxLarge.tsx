import '../Styles/TextBox.css';

import React from 'react';

/**
 * @interface LargeTextBoxProps
 * @property {string} value - The current text value of the textarea.
 * @property {(newValue: string) => void} onChange - Callback function for when the text value changes.
 * @property {() => void} onSave - Callback function for when the save button is clicked.
 * @property {string} [placeholder] - The placeholder text to display when the textarea is empty.
 * @property {boolean} [disabled=false] - If true, the textarea and button will be disabled.
 */
interface LargeTextBoxProps {
  value: string;
  onChange: (newValue: string) => void;
  onSave: () => void;
  placeholder?: string;
  disabled?: boolean;
}

/**
 * A large textarea component paired with a save button, designed for significant text input.
 * @param {LargeTextBoxProps} props The props for the component.
 * @returns {JSX.Element}
 */
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
