import "@/assets/styles/pages/DocumentPage.css";

import { DocumentHeaderControls, PDFView } from "../components";
import React, { useEffect, useMemo, useState } from "react";
import {
  useSetHeaderControls,
  useSetHeaderSubtitle,
  useSetHeaderTitle,
} from "@/components/Layouts/providers/HeaderProvider";

import { AppliedJob } from "@/features/applications/types";
import { CircleProgress } from "@/components/UI/loaders/CircleProgress";
import { DocumentType } from "@/features/documents/types";
import { ViewChangesModal } from "@/features/documents/components/ViewChangesModal";
import { useApplication } from "@/features/applications/providers/ApplicationProvider";
import { useDocumentManager } from "@/features/documents/hooks/useDocumentManager";

/**
 * Renders the main page for viewing, generating, and managing documents (resumes and cover letters)
 * associated with a selected job application. It handles different display states like loading,
 * generating, displaying the document, or showing an error.
 * @returns {React.FC} The DocumentPage component.
 */
export const DocumentPage: React.FC = () => {
  const { jobs, selectedJob, loading: appLoading } = useApplication();

  const [isChangesModalOpen, setIsChangesModalOpen] = useState(false);
  const [docType, setDocType] = useState<DocumentType>("resume");
  const [jobsWithDoc, setJobsWithDoc] = useState<AppliedJob[]>([]);
  const [jobsWithoutDoc, setJobsWithoutDoc] = useState<AppliedJob[]>([]);

  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  const {
    displayStatus,
    localPdfPath,
    localJsonData,
    generate,
    error,
    doesFileExist,
  } = useDocumentManager(
    selectedJob?.RoleID ?? null,
    selectedJob?.CompanyName ?? "",
    selectedJob?.JobTitle ?? "",
    docType
  );

  const isGenerating = displayStatus === 'generating';

  useEffect(() => {
    let isMounted = true;

    const sortJobs = async () => {
      const withDoc: AppliedJob[] = [];
      const withoutDoc: AppliedJob[] = [];

      for (const job of jobs) {
        const hasDoc = await doesFileExist(job.RoleID, docType, job.CompanyName, job.JobTitle);
        if (hasDoc && isMounted) {
          withDoc.push(job);
        } else if (isMounted) {
          withoutDoc.push(job);
        }
      }
      
      if (isMounted) {
        setJobsWithDoc(withDoc);
        setJobsWithoutDoc(withoutDoc);
      }
    };
    sortJobs();

    return () => {
      isMounted = false;
    };
  }, [jobs, docType, doesFileExist]);

  const headerControls = useMemo(() => (

    <DocumentHeaderControls
      selectedDocType={docType}
      onDocTypeChange={setDocType}
      isJobSelected={!!selectedJob}
      onCreate={generate}
      isGenerating={isGenerating}
      showViewChangesButton={displayStatus === 'present'}
      isCreateDisabled={isGenerating || !selectedJob}
      onViewChanges={() => setIsChangesModalOpen(true)}
      jobsWithDoc={jobsWithDoc}
      jobsWithoutDoc={jobsWithoutDoc}
    />
  ), [docType, selectedJob, generate, isGenerating, displayStatus, jobsWithDoc, jobsWithoutDoc]);

  useEffect(() => {
    if (selectedJob) {
      setHeaderTitle(selectedJob.CompanyProperName);
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