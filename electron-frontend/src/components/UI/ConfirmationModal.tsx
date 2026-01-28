import '@/assets/styles/Components/Layouts/ApplicationListRow.css';

import React from 'react';

interface ConfirmationModalProps {
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
  title?: string;
}

/**
 * A generic modal component used across the entire application 
 * to ask the user to confirm a destructive or important action.
 * * @param {string} title - Optional title for the modal (defaults to "Are you sure?")
 * @param {string} message - The specific warning message to display.
 * @param {() => void} onConfirm - Function to execute if the user clicks "Confirm".
 * @param {() => void} onCancel - Function to execute if the user clicks "Cancel".
 */
export const ConfirmationModal: React.FC<ConfirmationModalProps> = ({ 
  message, 
  onConfirm, 
  onCancel,
  title = "Are you sure?" 
}) => (
    <div className="confirmation-modal-overlay">
        <div className="confirmation-modal">
            <h2>{title}</h2>
            <p>{message}</p>
            <div className="confirmation-modal-buttons">
                <button onClick={onCancel} className="cancel-button">Cancel</button>
                <button onClick={onConfirm} className="confirm-button">Confirm</button>
            </div>
        </div>
    </div>
);