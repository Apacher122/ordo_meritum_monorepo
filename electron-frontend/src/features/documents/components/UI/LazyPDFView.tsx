import React from 'react';

const PDFView = React.lazy(() =>
  import('./PDFView.js').then(module => ({ default: module.PDFView }))
);
export const LazyPDFView = (props: React.ComponentProps<typeof PDFView>) => (
  <React.Suspense fallback={<div>Loading PDF Viewer...</div>}>
    <PDFView {...props} />
  </React.Suspense>
);