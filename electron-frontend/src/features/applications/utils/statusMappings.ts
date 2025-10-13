import { ApplicationStatus, BackendApplicationStatus } from "../types";

/**
 * @type BackendApplicationStatus
 * Represents the raw status strings as they are sent from the backend.
 */


/**
 * A mapping object to convert raw backend statuses into the normalized,
 * human-readable format used by the frontend.
 */
const statusMap: Record<BackendApplicationStatus, ApplicationStatus> = {
  REJECTED: 'Rejected',
  OFFERED: 'Offered',
  OPEN: 'Open',
  CLOSED: 'Closed',
  MOVED: 'Moved',
  NOT_APPLIED: 'Not applied',
  GHOSTED: 'Ghosted',
  INTERVIEWING: 'Interviewing',
};

const reverseStatusMap: { [key in ApplicationStatus]: BackendApplicationStatus } = {
  'Rejected': 'REJECTED',
  'Offered': 'OFFERED',
  'Open': 'OPEN',
  'Closed': 'CLOSED',
  'Moved': 'MOVED',
  'Not applied': 'NOT_APPLIED',
  'Ghosted': 'GHOSTED',
  'Interviewing': 'INTERVIEWING',
};


/**
 * Normalizes a status string from the backend into the frontend's format.
 * If the backend status is unknown, it defaults to "To Apply".
 * @param {string} backendStatus - The raw status string from the backend.
 * @returns {ApplicationStatus} The normalized, human-readable status.
 */
export const normalizeStatus = (
  backendStatus: string
): ApplicationStatus => {
    if (backendStatus in statusMap) {
    return statusMap[backendStatus as BackendApplicationStatus];
  }
    console.warn(`Unknown backend status received: "${backendStatus}"`);
  return "Not applied";
};

export const denormalizeStatus = (
    frontendStatus: ApplicationStatus
): BackendApplicationStatus => {
    return reverseStatusMap[frontendStatus];
};

export { BackendApplicationStatus };
