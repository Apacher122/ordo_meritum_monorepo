import { BackendApplicationStatus } from "../types";
import { apiRequest } from "@/shared/utils/requests";

export interface ApplicationUpdatePayload {
  status: BackendApplicationStatus | null;
  date: Date | null;
}


/**
 * Sends a PATCH request to the backend to update a job application.
 * @param roleId The ID of the job application to update.
 * @param payload An object containing the updated fields of the job application.
 * The object should contain one or both of the following properties: `status`, `date`.
 * If `status` is present, it should be a BackendApplicationStatus string.
 * If `date` is present, it should be a Date object.
 * @returns A promise that resolves with the response from the backend.
 */
export const updateApplication = (
  roleId: number,
  payload: ApplicationUpdatePayload
): Promise<Response> => {
  const params = new URLSearchParams({
    roleId: String(roleId),
  })
  return apiRequest(`api/secure/applications/update?${params.toString()}`, {
    method: "PATCH",
    body: payload,
  });
};
