import React from 'react';
import { usePublicKey } from '../../features/security/providers/PublicKeyProvider';

/**
 * Determines the status text and color based on the public key loading state.
 * @param {boolean} loading
 * @param {string | null} error
 * @param {string | null} publicKey
 * @returns {{color: string, text: string}}
 */
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

/**
 * Renders the application footer, which displays the current connection
 * status to the backend.
 * @returns {React.FC}
 */
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