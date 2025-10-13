

import React from 'react';
import { usePublicKey } from '../../features/security/providers/PublicKeyProvider';

export const Footer = () => {
  const { loading, publicKey, error } = usePublicKey();

  const statusColor = loading
    ? 'gray'
    : error
    ? 'red'
    : publicKey
    ? 'green'
    : 'yellow';

  const statusText = loading
    ? 'Connecting...'
    : error
    ? 'Connection Error'
    : publicKey
    ? 'Connected'
    : 'Not Found';

  return (
    <footer>
      <div className="status-indicator" style={{ backgroundColor: statusColor }} />
      <span>Status: {statusText}</span>
    </footer>
  );
};