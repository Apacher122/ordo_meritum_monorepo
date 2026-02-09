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
  SalaryRange?: string;
  Description?: string;
  YearsOfExp?: string;
  EducationLevel?: string;
  ApplicantCount?: string;
  PostAge?: string;
  Requirements?: string[];
  NiceToHaves?: string[];
  Tools?: string[];
  ProgrammingLanguages?: string[];
  FrameworksAndLibraries?: string[];
  Databases?: string[];
  CloudTechnologies?: string[];
  IndustryKeywords?: string[];
  SoftSkills?: string[];
  Certifications?: string[];
  CompanyCulture?: string;
  CompanyValues?: string;
}

export interface ApplicationData {
  hasSummary: boolean;
  hasCoverLetter: boolean;
  isVeteran: boolean;
  isDisabled: boolean;
  raceMentioned: boolean;
}