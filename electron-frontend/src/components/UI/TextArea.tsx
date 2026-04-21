import React from "react";

interface TextAreaProps extends Omit<React.TextareaHTMLAttributes<HTMLTextAreaElement>, 'value'> {
  value: string | undefined | null;
}

export const TextArea: React.FC<TextAreaProps> = ({ value, className = "textarea", ...props }) => {
  return (
    <textarea
      value={value || ""}
      className={className}
      {...props}
    />
  );
};