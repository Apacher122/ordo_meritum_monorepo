export interface ResumeChanges {
  type: "resume";
  summary: SummaryChanges[];
  experiences: ResumeChangeExperience[];
  skills: ResumeChangeSkill[];
  projects: ResumeChangeProject[];
}

export interface SummaryChanges {
  sentence: string;
  justification_for_change: string;
  is_new_suggestion: boolean;
}

export interface ResumeChangeSkill {
  category: string;
  skill: string[];
  justification_for_changes: string;
}

export interface ResumeChangeExperience {
  position: string;
  company: string;
  start: string;
  end: string;
  bulletPoints: ResumeChangeDescription[];
}

export interface ResumeChangeProject {
  name: string;
  role: string;
  status: string;
  bulletPoints: ResumeChangeDescription[];
}

export interface ResumeChangeDescription {
  text: string;
  justification_for_change: string;
  is_new_suggestion: boolean;
}
