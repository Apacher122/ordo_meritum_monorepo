import "@/assets/styles/Components/UI/DocumentHeaderControls.css";

import { AppliedJob } from "@/features/applications/types";
import { CircleProgress } from "@/components/UI/loaders/CircleProgress";
import { DocumentType } from "../../types";
import { DualJobSelectorDropdown } from "@/features/applications/components/DualJobSelectorDropDown";
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
  jobsWithDoc: AppliedJob[];
  jobsWithoutDoc: AppliedJob[];
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
  jobsWithDoc,
  jobsWithoutDoc,
}) => {
  return (
    <div className="document-header-controls">
      <DualJobSelectorDropdown
        jobsWithDoc={jobsWithDoc}
        jobsWithoutDoc={jobsWithoutDoc}
      />
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