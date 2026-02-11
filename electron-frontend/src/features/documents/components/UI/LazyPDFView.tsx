import { PDFViewProps } from './PDFView';
import React from 'react';

const PDFView = React.lazy<React.ComponentType<PDFViewProps>>(() =>
  import('./PDFView').then(module => ({ default: module.PDFView }))
);

/**
 * A lazy-loaded version of the PDFView component.
 * It displays a "Loading PDF Viewer..." message until the PDFView component is loaded.
 * @param {PDFViewProps} props - The props to pass to the PDFView component.
 * @returns {React.Component} - The lazy-loaded PDFView component.
 */
export const LazyPDFView = (props: PDFViewProps) => (
  <React.Suspense fallback={<div>Loading PDF Viewer...</div>}>
    <PDFView {...props} />
  </React.Suspense>
);