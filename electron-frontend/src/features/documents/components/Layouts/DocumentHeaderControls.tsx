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
  jobsNotAppliedNoDoc: AppliedJob[];
}


/**
 * A component that displays the header controls for the document viewer.
 * It includes dropdown selectors for all applications and ready to analyze applications.
 * It also includes buttons for regenerating the document and viewing changes.
 * @param {DocumentHeaderControlsProps} props - The component props.
 * @param {DocumentType} selectedDocType - The selected document type (resume or cover-letter).
 * @param {(docType: DocumentType) => void} onDocTypeChange - The function to call when the document type changes.
 * @param {boolean} isJobSelected - Whether a job is selected or not.
 * @param {() => void} onCreate - The function to call when the create button is clicked.
 * @param {() => void} onViewChanges - The function to call when the view changes button is clicked.
 * @param {boolean} showViewChangesButton - Whether to show the view changes button or not.
 * @param {boolean} isGenerating - Whether the document is generating or not.
 * @param {boolean} isCreateDisabled - Whether the create button is disabled or not.
 * @param {AppliedJob[]} jobsWithDoc - The array of job applications with documents.
 * @param {AppliedJob[]} jobsWithoutDoc - The array of job applications without documents.
 * @param {AppliedJob[]} jobsNotAppliedNoDoc - The array of job applications that are not applied and do not have a document.
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
  jobsNotAppliedNoDoc,
}) => {
  return (
    <div className="document-header-controls">
      <div className="dropdowns-row">
        <div className="dropdown-field">
          <label htmlFor="all-jobs-selector" className="dropdown-label">All Applications</label>
          <DualJobSelectorDropdown
            id="all-jobs-selector"
            jobsWithDoc={jobsWithDoc}
            jobsWithoutDoc={jobsWithoutDoc}
          />
        </div>

        <div className="dropdown-field">
          <label htmlFor="not-applied-selector" className="dropdown-label">Ready to Analyze</label>
          <DualJobSelectorDropdown
            id="not-applied-selector"
            jobsWithDoc={[]}
            jobsWithoutDoc={jobsNotAppliedNoDoc}
          />
        </div>
      </div>

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