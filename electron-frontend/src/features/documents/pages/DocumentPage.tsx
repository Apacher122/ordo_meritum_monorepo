import "@/assets/styles/pages/DocumentPage.css";

import { DocumentHeaderControls, LazyPDFView } from "../components";
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
import { useLocation } from "react-router-dom";

/**
 * A page that displays the document viewer for a job application.
 * It takes the job ID, company name, job title, and document type as parameters.
 * It fetches the document data from the server and displays the document in a PDF viewer.
 * It also provides a button to generate the document and a button to view changes.
 * If the document does not exist, it displays a message indicating that the document can be created.
 * If the document is being generated, it displays a loading indicator.
 * If the document generation fails, it displays an error message.
 * If the document exists, it displays the document in a PDF viewer.
 */
export const DocumentPage: React.FC = () => {
  const { jobs, selectedJob, loading: appLoading } = useApplication();
  const location = useLocation();

  const [isChangesModalOpen, setIsChangesModalOpen] = useState(false);
  
  const [docType, setDocType] = useState<DocumentType>(
    (location.state?.initialDocType as DocumentType) || "resume"
  );
  
  const [jobsWithDoc, setJobsWithDoc] = useState<AppliedJob[]>([]);
  const [jobsWithoutDoc, setJobsWithoutDoc] = useState<AppliedJob[]>([]);
  const [jobsNotAppliedNoDoc, setJobsNotAppliedNoDoc] = useState<AppliedJob[]>([]);

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
      const notAppliedNoDoc: AppliedJob[] = [];

      for (const job of jobs) {
        const hasDoc = await doesFileExist(job.RoleID, docType, job.CompanyName, job.JobTitle);
        if (hasDoc && isMounted) {
          withDoc.push(job);
        } else if (isMounted) {
          withoutDoc.push(job);
          if (job.ApplicationStatus === "Not applied") {
            notAppliedNoDoc.push(job);
          }
        }
      }
      
      if (isMounted) {
        setJobsWithDoc(withDoc);
        setJobsWithoutDoc(withoutDoc);
        setJobsNotAppliedNoDoc(notAppliedNoDoc);
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
      jobsNotAppliedNoDoc={jobsNotAppliedNoDoc}
    />
  ), [
    docType,
    selectedJob,
    generate,
    isGenerating,
    displayStatus,
    jobsWithDoc,
    jobsWithoutDoc,
    jobsNotAppliedNoDoc
  ]);

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
        { const temp = `static://${localPdfPath}`;
        return (
          <div className="pdf-container">
            <LazyPDFView file={temp} />
          </div>
        ); }

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