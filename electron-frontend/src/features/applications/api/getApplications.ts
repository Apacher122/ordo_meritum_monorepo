import { AppliedJob } from "../types";
import { apiRequest } from "@/shared/utils/requests";

/**
 * Fetches the list of job applications for the authenticated user.
 * @param {string} token - The Bearer token for API requests.
 * @returns {Promise<AppliedJob[]>} A promise that resolves to the list of jobs.
 */
export const getApplications = (token: string): Promise<AppliedJob[]> => {
  return apiRequest<AppliedJob[]>("api/auth/apps/track/list", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
};