import { z } from "zod";

export const CompilationResultSchema = z.object({
  user_id: z.string(),
  job_id: z.number().int(),
  success: z.boolean(),
  document_type: z.string().optional(),
  download_url: z.string().optional(),
  changes_url: z.string().optional(),
  error: z.string().optional(),
});

export type CompilationResult = z.infer<typeof CompilationResultSchema>;
