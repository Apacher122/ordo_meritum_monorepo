import { apiRequest } from '@/shared/utils/requests';

/**
 * Deletes an application from the database.
 * @param jobId The ID of the job application to delete.
 */
export const deleteApplication = async (jobId: number): Promise<void> => {
  
  const params = new URLSearchParams({
    jobId: String(jobId),
  })
  await apiRequest<void>(`api/secure/applications/delete?${params.toString()}`, {
    method: 'DELETE',
  });
};
