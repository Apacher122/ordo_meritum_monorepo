import '@/assets/styles/App.css';

import { AppRoutes } from './router';
import { pdfjs } from 'react-pdf';

pdfjs.GlobalWorkerOptions.workerSrc = './pdf.worker.min.mjs';

function App() {
  return <AppRoutes />;
}

export default App;