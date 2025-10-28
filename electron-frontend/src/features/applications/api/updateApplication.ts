import { BackendApplicationStatus } from "../types";
import { apiRequest } from "@/shared/utils/requests";

export interface ApplicationUpdatePayload {
  job_id: number;
  job_title: string | null;
  website: string | null;
  application_status: BackendApplicationStatus | null;
  interview_count: number | null;
  initial_application_date: Date | null;
}

/**
 * Sends a PATCH request to the backend to update a job application.
 * @param {string} token The Bearer token for API requests.
 * @param {ApplicationUpdatePayload} payload The request body containing the updated job application information.
 * @returns {Promise<Response>} A promise resolving to the response from the backend.
 */
export const updateApplication = (
  token: string,
  payload: ApplicationUpdatePayload
): Promise<Response> => {
  const body = {
    payload: payload,
  };
  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
  return apiRequest(`api/auth/apps/update`, {
    method: "PATCH",
    headers: headers,
    body: body,
  });
};
