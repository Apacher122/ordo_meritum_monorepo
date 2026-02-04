import { apiRequest } from '@/shared/utils/requests';

/**
 * Deletes an application from the database.
 * @param jobId The ID of the job application to delete.
 */
export const deleteApplication = async (jobId: number): Promise<void> => {
  await apiRequest<void>(`api/auth/apps/delete`, {
    method: 'DELETE',
    body: {payload: {job_id: jobId}},
  });
};
