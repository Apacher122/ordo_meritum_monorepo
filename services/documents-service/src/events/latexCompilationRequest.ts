import { z } from "zod";

export const DocumentRequestSchema = z.object({
  apiKey: z.string(),
  payload: z.any(),
  options: z.object({
    jobId: z.preprocess((value) => Number(value), z.number().int().positive()),
    docType: z.enum(['resume', 'cover-letter']),
    llm: z.enum(['openai', 'cohere', 'ollama', 'gemini', 'groq', 'claude', 'cerebras']).optional(),
    getNew: z.preprocess((value) => value === "true" || value === true, z.boolean()).optional(),
    corrections: z.array(z.string()).optional().default([]),
  })
})

export type DocumentRequest = z.infer<typeof DocumentRequestSchema>;

export const UserInfoSchema = z.object({
  first_name: z.string(),
  last_name: z.string(),
  current_location: z.string(),
  email: z.string().email(),
  github: z.string(),
  linkedin: z.string(),
  mobile: z.string(),
  summary: z.string(),
});

export const EducationPayloadSchema = z.object({
  coursework: z.string(),
  degree: z.string(),
  location: z.string(),
  school: z.string(),
  start_end: z.string(),
});

export const SummaryPayloadSchema = z.object({
  sentence: z.string(),
  justification_for_change: z.string(),
  is_new_suggestion: z.boolean(),
});

export const SkillPayloadSchema = z.object({
  category: z.string(),
  justification_for_changes: z.string(),
  skill: z.array(z.string()),
});

const BulletPointSchema = z.object({
  text: z.string(),
  is_new_suggestion: z.boolean(),
  justification_for_change: z.string(),
});

export const ExperiencePayloadSchema = z.object({
  position: z.string(),
  company: z.string(),
  start: z.string(), 
  end: z.string(),
  bulletPoints: z.array(BulletPointSchema),
});

export const ProjectPayloadSchema = z.object({
  name: z.string(),
  role: z.string(),
  bulletPoints: z.array(BulletPointSchema),
});

export const ResumePayloadSchema = z.object({
  summary: z.array(SummaryPayloadSchema),
  skills: z.array(SkillPayloadSchema),
  experiences: z.array(ExperiencePayloadSchema),
  projects: z.array(ProjectPayloadSchema).optional().nullable(),
});

export const CoverLetterBody = z.object({
  about: z.string(),
  experience: z.string(),
  whatIBring: z.string(),
  revisionSummary: z.string(),
})

export const CoverLetterPayloadSchema = z.object({
  companyProperName: z.string(),
  jobTitle: z.string(),
  body: CoverLetterBody,
})


export const CompilationRequestSchema = z.object({
  jobID: z.number().int(),
  userID: z.string(),
  companyName: z.string(),
  docType: z.string(),
  userInfo: UserInfoSchema,
  educationInfo: EducationPayloadSchema,
  resume: ResumePayloadSchema,
  coverLetter: CoverLetterPayloadSchema.optional().nullable(),
});

export type CompilationRequest = z.infer<typeof CompilationRequestSchema>;
export type UserInfo = z.infer<typeof UserInfoSchema>;
export type EducationInfo = z.infer<typeof EducationPayloadSchema>;
export type ResumePayload = z.infer<typeof ResumePayloadSchema>;
export type CoverLetterPayload = z.infer<typeof CoverLetterPayloadSchema>;
export type CoverLetterBodyType = z.infer<typeof CoverLetterBody>;