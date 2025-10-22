import React, { ReactNode, createContext, useContext, useMemo } from "react";

import { useAuth } from "../../auth/providers/AuthProvider";
import { usePublicKeyStream } from "../hooks/usePublicKeyStream";

interface PublicKeyContextType {
  publicKey: string | null;
  loading: boolean;
  error: string | null;
}

const PublicKeyContext = createContext<PublicKeyContextType | undefined>(undefined);

export const PublicKeyProvider = ({ children }: { children: ReactNode }) => {
  const { user } = useAuth();   const { publicKey, loading, error } = usePublicKeyStream(user);   
  const value = useMemo(() => ({ publicKey, loading, error }), [publicKey, loading, error]);
  return (
    <PublicKeyContext.Provider value={value}>
      {children}
    </PublicKeyContext.Provider>
  );
};

export const usePublicKey = () => {
  const context = useContext(PublicKeyContext);
  if (!context) {
    throw new Error("usePublicKey must be used within a PublicKeyProvider");
  }
  return context;
};