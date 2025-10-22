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
