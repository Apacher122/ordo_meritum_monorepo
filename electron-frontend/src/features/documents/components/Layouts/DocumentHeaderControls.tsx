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

/**
 * DocumentHeaderControls is a component that displays a header for the document page
 * 
 * @param {DocumentHeaderControlsProps} props - The props object
 * @param {DocumentType} props.selectedDocType - The selected document type
 * @param {(docType: DocumentType) => void} props.onDocTypeChange - The callback function when the document type changes
 * @param {boolean} props.isJobSelected - Whether a job is selected or not
 * @param {() => void} props.onCreate - The callback function when the user clicks on the create button
 * @param {() => void} props.onViewChanges - The callback function when the user clicks on the view changes button
 * @param {boolean} props.showViewChangesButton - Whether to show the view changes button or not
 * @param {boolean} props.isGenerating - Whether the document is being generated or not
 * @param {boolean} props.isCreateDisabled - Whether the create button is disabled or not
 * @param {AppliedJob[]} props.jobsWithDoc - The list of jobs with a document
 * @param {AppliedJob[]} props.jobsWithoutDoc - The list of jobs without a document
 */
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