import { apiRequest } from '@/shared/utils/requests';

/**
 * Deletes an application from the database.
 * @param jobId The ID of the job application to delete.
 */
export const deleteApplication = async (
  token: string,
  jobId: number
): Promise<void> => {
  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
  const params = new URLSearchParams({
    jobId: String(jobId),
  })
  await apiRequest<void>(`api/auth/apps/delete?${params.toString()}`, {
    method: 'DELETE',
    headers: headers,
    body: {payload: {job_id: jobId}},
  });
};
