import "react-pdf/dist/esm/Page/AnnotationLayer.css";
import "react-pdf/dist/esm/Page/TextLayer.css";
import "@/assets/styles/Components/PdfView.css";

import { Document, Page, pdfjs } from "react-pdf";
import React, { useState } from "react";

pdfjs.GlobalWorkerOptions.workerSrc = "./pdf.worker.min.mjs";

interface PDFViewProps {
  file: Blob | string | null;
}

interface PDFViewProps {
  file: Blob | string | null;
}

export const PDFView: React.FC<PDFViewProps> = ({ file }) => {
  const [numPages, setNumPages] = useState<number>(0);
  const [pageNumber, setPageNumber] = useState<number>(1);
  const [scale, setScale] = useState<number>(1);

  const onDocumentLoadSuccess = ({ numPages }: { numPages: number }) => {
    setNumPages(numPages);
    setPageNumber(1);
  };

  const goToPrevPage = () =>
    setPageNumber((prevPageNumber) => Math.max(prevPageNumber - 1, 1));

  const goToNextPage = () =>
    setPageNumber((prevPageNumber) => Math.min(prevPageNumber + 1, numPages));

  const zoomIn = () => setScale(prevScale => Math.min(prevScale + 0.1, 2));
  const zoomOut = () => setScale(prevScale => Math.max(prevScale - 0.1, 0.5));
  const resetZoom = () => setScale(1);

  if (!file) {
    return (
      <div className="pdf-view-container centered-feedback">
        <p>No document to display.</p>
      </div>
    );
  }

  return (
    <div className="pdf-view-container">
      <div className="pdf-toolbar">
        <div className="pagination-controls">
          <button onClick={goToPrevPage} disabled={pageNumber <= 1}>
            Prev
          </button>
          <span>
            Page {pageNumber} of {numPages}
          </span>
          <button onClick={goToNextPage} disabled={pageNumber >= numPages}>
            Next
          </button>
        </div>
        <div className="zoom-controls">
          <button onClick={zoomOut} disabled={scale <= 0.5}>-</button>
          <button 
            onClick={resetZoom} 
            title="Reset Zoom" 
            className="zoom-level"
          >
            {Math.round(scale * 100)}%
          </button>
          <button onClick={zoomIn} disabled={scale >= 2}>+</button>
        </div>
      </div>

      <div className="pdf-document-container">
        <Document
          file={file}
          onLoadSuccess={onDocumentLoadSuccess}
          loading={<div className="centered-feedback">Loading PDF...</div>}
          error={<div className="centered-feedback error-message">Failed to load PDF.</div>}
        >
          <div className="pdf-page-wrapper">
            <Page
              pageNumber={pageNumber}
              scale={scale}
              renderTextLayer={true} 
              className="pdf-page-dark"
            />
          </div>
        </Document>
      </div>
    </div>
  );
};