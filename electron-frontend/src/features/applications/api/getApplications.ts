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