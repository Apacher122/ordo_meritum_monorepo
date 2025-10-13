import { LlmProvider } from "@/shared/types/index.js";

export type AssignableFeature =
  | "matchSummary"
  | "resumeGeneration"
  | "coverLetterGeneration";

export interface Settings {
  apiKeys: Partial<Record<LlmProvider, string>>;
  featureAssignments: Record<AssignableFeature, LlmProvider>;
}
