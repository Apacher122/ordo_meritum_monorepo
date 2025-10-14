import "@/assets/styles/Components/UI/DocumentHeaderControls.css";

import { CircleProgress } from "@/components/UI/loaders/CircleProgress";
import { DocumentType } from "../../types";
import { JobSelectorDropdown } from "@/features/applications/components/JobSelectorDropdown";
import React from "react";

interface DocumentHeaderControlsProps {
  selectedDocType: DocumentType;
  onDocTypeChange: (docType: DocumentType) => void;
  isJobSelected: boolean;
  onCreate: () => void;
  onViewChanges: () => void;
  showViewChangesButton: boolean;
  isGenerating: boolean;
  isCreateDisabled: boolean;
}

export const DocumentHeaderControls: React.FC<DocumentHeaderControlsProps> = ({
  selectedDocType,
  onDocTypeChange,
  isJobSelected,
  onCreate,
  onViewChanges,
  showViewChangesButton,
  isGenerating,
  isCreateDisabled,
}) => {
  return (
    <div className="document-header-controls">
      <JobSelectorDropdown />
      {isJobSelected && (
        <>
          <div className="doc-type-tabs">
            <button
              className={`tab-button ${selectedDocType === "resume" ? "active" : ""}`}
              onClick={() => onDocTypeChange("resume")}
              disabled={isGenerating}
            >
              Resume
            </button>
            <button
              className={`tab-button ${selectedDocType === "cover-letter" ? "active" : ""}`}
              onClick={() => onDocTypeChange("cover-letter")}
              disabled={isGenerating}
            >
              Cover Letter
            </button>
          </div>
          
          <div className="header-action-buttons">
            {showViewChangesButton && (
              <button
                className="button secondary"
                onClick={onViewChanges}
                disabled={isGenerating}
              >
                View Changes
              </button>
            )}
            <button
              className="button"
              onClick={onCreate}
              disabled={isCreateDisabled || isGenerating}
            >
              {isGenerating ? "Processing..." : `Regenerate`}
            </button>
            {isGenerating && <CircleProgress size={24} strokeWidth={3} />}
          </div>
        </>
      )}
    </div>
  );
};