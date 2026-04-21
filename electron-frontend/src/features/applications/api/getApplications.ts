import { AppliedJob } from "../types";
import { apiRequest } from "@/shared/utils/requests";

/**
 * Sends the user's Firebase UID to the backend to register or log in the user
 * in the application's own database.
 * @returns {Promise<void>} A promise that resolves when the request is complete.
 */
export const getApplications = (): Promise<AppliedJob[]> => {
  return apiRequest<AppliedJob[]>("api/auth/apps/track/list", {
    method: "GET",
  });
};

/**
 * Retrieves a tracked job posting by its ID.
 * @param {number} job_id The ID of the job posting to retrieve.
 * @returns {Promise<AppliedJob>} A promise that resolves to the retrieved job posting if successful, or rejects with an error if failed.
 */
export const getApplicationById = (job_id: number): Promise<AppliedJob> => {
  return apiRequest<AppliedJob>(`api/auth/apps/${job_id}`, {
    method: "GET",
  });
}