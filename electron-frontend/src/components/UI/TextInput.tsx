import React from "react";

interface TextInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'value'> {
  value: string | undefined | null;
}

export const TextInput: React.FC<TextInputProps> = ({ value, className = "input", ...props }) => {
  return (
    <input
      value={value || ""}
      className={className}
      {...props}
    />
  );
};