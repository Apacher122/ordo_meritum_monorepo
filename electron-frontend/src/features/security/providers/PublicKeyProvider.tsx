import React, { ReactNode, createContext, useContext, useMemo } from "react";

import { useAuth } from "../../auth/providers/AuthProvider";
import { usePublicKeyStream } from "../hooks/usePublicKeyStream";

interface PublicKeyContextType {
  publicKey: string | null;
  loading: boolean;
  error: string | null;
}

const PublicKeyContext = createContext<PublicKeyContextType | undefined>(undefined);

/**
 * A React Context Provider that provides the public key for the currently authenticated user.
 * It uses the useAuth hook to get the user and the usePublicKeyStream hook to get the public key.
 * The Context Provider provides the public key, loading state, and any error that might have occurred.
 * The public key is only available if a user is logged in.
 */
export const PublicKeyProvider = ({ children }: { children: ReactNode }) => {
  const { user } = useAuth();   const { publicKey, loading, error } = usePublicKeyStream(user);   
  const value = useMemo(() => ({ publicKey, loading, error }), [publicKey, loading, error]);
  return (
    <PublicKeyContext.Provider value={value}>
      {children}
    </PublicKeyContext.Provider>
  );
};

/**
 * Retrieves the public key for the currently authenticated user.
 * It throws an error if the context is not available, i.e. when it is not used within a PublicKeyProvider.
 * The public key is only available if a user is logged in.
 * @returns {PublicKeyContextType} An object containing the public key, loading state, and any error that might have occurred.
 */
export const usePublicKey = () => {
  const context = useContext(PublicKeyContext);
  if (!context) {
    throw new Error("usePublicKey must be used within a PublicKeyProvider");
  }
  return context;
};