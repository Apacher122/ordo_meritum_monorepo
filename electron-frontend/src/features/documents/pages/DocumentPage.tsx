import "@/assets/styles/pages/DocumentPage.css";
import React, { useEffect, useMemo, useState } from "react";
import { DocumentHeaderControls } from "../components";
import { DocumentType } from "../types";
import { PDFView } from "../components";
import { useApplication } from "../../applications/providers/ApplicationProvider";
import { useDocumentManager } from "../hooks/useDocumentManager";
import { CircleProgress } from "@/components/UI/loaders/CircleProgress";
import {
  useSetHeaderTitle,
  useSetHeaderSubtitle,
  useSetHeaderControls,
} from "@/components/Layouts/providers/HeaderProvider";
import { ViewChangesModal } from "../components/ViewChangesModal";

export const DocumentPage: React.FC = () => {
  const { selectedJob, loading: appLoading } = useApplication();
  const [isChangesModalOpen, setIsChangesModalOpen] = useState(false);
  const [docType, setDocType] = useState<DocumentType>();

  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  const {
    displayStatus,
    localPdfPath,
    localJsonData,
    generate,
    error,
  } = useDocumentManager(
    selectedJob?.RoleID ?? null,
    selectedJob?.CompanyName ?? "",
    selectedJob?.JobTitle ?? "",
    docType ?? "resume"
  );

  const isGenerating = displayStatus === 'generating';

  const headerControls = useMemo(() => (
    <DocumentHeaderControls
      selectedDocType={docType ?? "resume"}
      onDocTypeChange={setDocType}
      isJobSelected={!!selectedJob}
      onCreate={generate}
      isGenerating={isGenerating}
      showViewChangesButton={displayStatus === 'present'}
      isCreateDisabled={isGenerating} 
      onViewChanges={() => setIsChangesModalOpen(true)}
    />
  ), [docType, selectedJob, generate, isGenerating, displayStatus]);

  useEffect(() => {
    if (selectedJob) {
      setHeaderTitle(selectedJob.CompanyName);
      setHeaderSubtitle(selectedJob.JobTitle);
    } else {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
    }
    setHeaderControls(headerControls);

    return () => {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
      setHeaderControls(null);
    };
  }, [selectedJob, headerControls, setHeaderTitle, setHeaderSubtitle, setHeaderControls]);

  if (appLoading) {
    return <div>Loading Application...</div>;
  }

  if (!selectedJob) {
    return <div className="page-content-placeholder">Please select a job application to view documents.</div>;
  }

  const renderContent = () => {
    switch (displayStatus) {
      case 'checking':
        return <div className="page-content-placeholder"><CircleProgress /></div>;

      case 'generating':
        return (
          <div className="page-content-placeholder">
            <CircleProgress />
            <p>Generating {docType?.replace("-", " ")}... This may take a moment.</p>
            <p className="subtle-warning">You can safely navigate away from this page.</p>
          </div>
        );

      case 'present':
        return (
          <div className="pdf-container">
            <PDFView file={localPdfPath} />
          </div>
        );

      case 'failed':
        return (
          <div className="page-content-placeholder error-message">
            <h2>Generation Failed</h2>
            <p>{error || "An unknown error occurred."}</p>
          </div>
        );

      case 'idle':
      default:
        return (
          <div className="page-content-placeholder">
            <h2>No {docType?.replace("-", " ")} Exists</h2>
            <p>You can create one using the button in the header.</p>
          </div>
        );
    }
  };

  return (
    <div className="document-page">
      {renderContent()}
      {localJsonData && (
        <ViewChangesModal
          isOpen={isChangesModalOpen}
          onClose={() => setIsChangesModalOpen(false)}
          changes={localJsonData}
        />
      )}
    </div>
  );
};

