import { ApplicationStatus } from "./statuses";

export interface AppliedJob {
  RoleID: number;
  JobTitle: string;
  CompanyName: string;
  CompanyProperName: string;
  Website: string;
  ApplicationStatus: ApplicationStatus;
  UserApplied: boolean;
  InterviewCount: number;
  InitialApplicationDate: Date;
}