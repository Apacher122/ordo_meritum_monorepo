import "react-pdf/dist/Page/AnnotationLayer.css";
import "react-pdf/dist/Page/TextLayer.css";
import "@/assets/styles/Components/PdfView.css";

import { Document, Page, pdfjs } from "react-pdf";
import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";

pdfjs.GlobalWorkerOptions.workerSrc = "/pdf.worker.min.mjs";

export interface PDFViewProps {
  file: Blob | string | null;
}

/**
 * A React component for displaying a PDF document.
 * It will load the PDF document via IPC if the file is a local path.
 * It will display the PDF document with pagination controls and zoom controls.
 * It will display a loading message while the PDF document is being loaded.
 * It will display an error message if the PDF document fails to load.
 * @param {{ file: Blob|string|null }} props
 * @returns {JSX.Element}
 */
export const PDFView: React.FC<PDFViewProps> = ({ file }) => {
  const [numPages, setNumPages] = useState<number>(0);
  const [pageNumber, setPageNumber] = useState<number>(1);
  const [scale, setScale] = useState<number>(1);
  const [pdfBytes, setPdfBytes] = useState<Uint8Array | null>(null);
  const [pdfLoaded, setPdfLoaded] = useState<boolean>(false);
  const isMounted = useRef(true);

  useEffect(() => {
    isMounted.current = true;
    return () => { isMounted.current = false; };
  }, []);

  useEffect(() => {
    const loadPdfData = async () => {
      if (!file) {
        setPdfBytes(null);
        return;
      }
      const isLocalPath = typeof file === "string" && (file.startsWith("static://") || file.includes(".pdf"));
      if (isLocalPath && window.appAPI?.files?.readPdfBytes) {
        try {
          const bytes = await window.appAPI.files.readPdfBytes(file);
          if (isMounted.current) setPdfBytes(bytes);
        } catch (err) {
          console.error("PDFView: IPC Load failed", err);
        }
      }
    };
    loadPdfData();
  }, [file]);

  const memoizedFile = useMemo(() => {
    if (pdfBytes) return { data: pdfBytes };
    return file;
  }, [file, pdfBytes]);

  const onDocumentLoadSuccess = useCallback(({ numPages }: { numPages: number }) => {
    if (isMounted.current) {
      setNumPages(numPages);
      setPageNumber(1);
      setPdfLoaded(true);
    }
  }, []);

  const onDocumentLoadError = useCallback(() => {
    if (isMounted.current) setPdfLoaded(false);
  }, []);

  const goToPrevPage = () => setPageNumber(prev => Math.max(prev - 1, 1));
  const goToNextPage = () => setPageNumber(prev => Math.min(prev + 1, numPages));
  const zoomIn = () => setScale(prev => Math.min(prev + 0.1, 2));
  const zoomOut = () => setScale(prev => Math.max(prev - 1, 0.5));

  if (!file) return <div className="centered-feedback">No document to display.</div>;

  return (
    <div className="pdf-view-container">
      <div className="pdf-toolbar">
        <div className="pagination-controls">
          <button onClick={goToPrevPage} disabled={pageNumber <= 1}>Prev</button>
          <span>Page {pageNumber} of {numPages}</span>
          <button onClick={goToNextPage} disabled={pageNumber >= numPages}>Next</button>
        </div>
        <div className="zoom-controls">
          <button onClick={zoomOut} disabled={scale <= 0.5}>-</button>
          <span className="zoom-level">{Math.round(scale * 100)}%</span>
          <button onClick={zoomIn} disabled={scale >= 2}>+</button>
        </div>
      </div>

      <div className="pdf-document-container">
        <Document
          key={typeof file === 'string' ? file : 'blob-source'}
          file={memoizedFile}
          onLoadSuccess={onDocumentLoadSuccess}
          onLoadError={onDocumentLoadError}
          loading={<div className="centered-feedback">Loading PDF...</div>}
        >
          {pdfLoaded && (
            <div className="pdf-page-wrapper">
              <Page
                pageNumber={pageNumber}
                scale={scale}
                renderTextLayer={true} 
                renderAnnotationLayer={true}
                className="pdf-page-dark"
              />
            </div>
          )}
        </Document>
      </div>
    </div>
  );
};