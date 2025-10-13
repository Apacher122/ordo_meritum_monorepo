export interface UserInfo {
  first_name?: string;
  last_name?: string;
  summary?: string;
  location?: string;
  current_location?: string;
  mobile?: string;
  email?: string;
  github?: string;
  linkedin?: string;
}
export interface EducationInfo {
  school?: string;
  degree?: string;
  start_end?: string;
  location?: string;
  coursework?: string;
  undergraduate_coursework?: string;
  graduate_coursework?: string;
  education_summary?: string;
}

export interface Experience {
  id: string; 
  company?: string;
  jobTitle?: string;
  years?: string;
  bulletPoints?: {text: string}[];
}

export interface Project {
  id: string; 
  name?: string;
  description?: string;
  years?: string;
  bulletPoints?: {text: string}[];
}

export interface ResumeData {
  skills?: {skill: string}[];
  experiences?: Experience[];
  projects?: Project[];
}

export interface CoverLetterData {
  about?: string;
  experience?: string;
  whatIBring?: string;
}

export interface AboutMeData {
  essay?: string;
}

export interface WritingSample {
  fileName: string;
  content: string;
}

export interface UserProfile {
  userInfo: UserInfo;
  education: EducationInfo;
  resume: ResumeData;
  coverLetter: CoverLetterData;
  aboutMe: AboutMeData;
  writingSamples?: WritingSample[];
}