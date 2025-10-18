import React from 'react';
import { usePublicKey } from '../../features/security/providers/PublicKeyProvider';

const getStatus = (loading: boolean, error: string | null, publicKey: string | null) => {
  if (loading) {
    return { color: "gray", text: "Connecting..." };
  }
  if (error) {
    return { color: "red", text: "Connection Error" };
  }
  if (publicKey) {
    return { color: "green", text: "Connected" };
  }
  return { color: "yellow", text: "Not Found" };
};

export const Footer = () => {
  const { loading, publicKey, error } = usePublicKey();

const { color: statusColor, text: statusText } = getStatus(loading, error, publicKey);

  return (
    <footer>
      <div className="status-indicator" style={{ backgroundColor: statusColor }} />
      <span>Status: {statusText}</span>
    </footer>
  );
};