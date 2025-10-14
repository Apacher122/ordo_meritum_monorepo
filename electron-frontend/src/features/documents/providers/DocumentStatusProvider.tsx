import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';

import { DocumentType } from '../types';
import { useAuth } from '@/features/auth/providers/AuthProvider';
import { webSocketService } from '@/shared/api/WebSocketService';

export interface DocumentStatus {
  status: 'PENDING' | 'COMPLETED' | 'FAILED';
  downloadUrl: string;
  changesUrl: string;
  error?: string;
}

export type JobDocumentsStatus = Map<DocumentType, DocumentStatus>;
type DocumentStatusState = Map<string, JobDocumentsStatus>;

interface DocumentStatusUpdateMessage {
  user_id: string;
  job_id: string;
  success: boolean;
  document_type: DocumentType;
  download_url: string;
  changes_url: string;
  error?: string;
}

interface DocumentStatusContextType {
  documentStatuses: DocumentStatusState;
  addPendingDocument: (jobId: string, docType: DocumentType) => void;
}

const DocumentStatusContext = createContext<DocumentStatusContextType | undefined>(undefined);

export const DocumentStatusProvider = ({ children }: { children: React.ReactNode }) => {
  const { user } = useAuth();
  const [documentStatuses, setDocumentStatuses] = useState<DocumentStatusState>(new Map());

  useEffect(() => {
    if (user) {
      const setupWebSocket = async () => {
        try {
          const token = await user.getIdToken();
          webSocketService.connect(user.uid, token); 
        } catch (error) {
          console.error("Failed to get Firebase ID token:", error);
          webSocketService.disconnect();
        }
      };
      
      setupWebSocket();

      webSocketService.onMessage((data: DocumentStatusUpdateMessage) => {
        console.log('Received document status update:', data);
        setDocumentStatuses((prev) => {
          const newStatuses = new Map(prev);
          const currentJob = new Map(newStatuses.get(String(data.job_id)) ?? []);
          currentJob.set(data.document_type, {
            status: data.success ? 'COMPLETED' : 'FAILED',
            downloadUrl: data.download_url,
            changesUrl: data.changes_url,
            error: data.error,
          });

          newStatuses.set(String(data.job_id), currentJob);
          return newStatuses;
        });
      });
    } else {
      webSocketService.disconnect();
    }

    return () => {
      webSocketService.disconnect();
    };
  }, [user]);

  const addPendingDocument = useCallback((jobId: string, docType: DocumentType) => {
    setDocumentStatuses((prev) => {
      const newStatuses = new Map(prev);
      const currentJob = new Map(newStatuses.get(jobId));
      currentJob.set(docType, { status: 'PENDING', downloadUrl: '', changesUrl: '', error: '' });
      newStatuses.set(jobId, currentJob);
      return newStatuses;
    });
  }, []);

  const value = { documentStatuses, addPendingDocument };

  return (
    <DocumentStatusContext.Provider value={value}>
      {children}
    </DocumentStatusContext.Provider>
  );
};

export const useDocumentStatus = () => {
  const context = useContext(DocumentStatusContext);
  if (context === undefined) {
    throw new Error('useDocumentStatus must be used within a DocumentStatusProvider');
  }
  return context;
};
