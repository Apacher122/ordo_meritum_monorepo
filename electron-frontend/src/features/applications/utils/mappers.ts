import { AppliedJob } from "../types";
import { normalizeStatus } from "./statusMappings";

/**
 * Transforms raw job application data from the backend into the frontend's expected format.
 * Adjusts timezone offsets and normalizes status strings.
 */
export const transformJobData = (job: AppliedJob): AppliedJob => {
  const dateFromBackend = new Date(job.InitialApplicationDate);
  const timezoneOffset = dateFromBackend.getTimezoneOffset() * 60000;
  const correctedDate = new Date(dateFromBackend.getTime() + timezoneOffset);

  return {
    ...job,
    ApplicationStatus: normalizeStatus(job.ApplicationStatus),
    InitialApplicationDate: correctedDate,
  };
};