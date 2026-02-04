import { LlmProvider } from "@/shared/types/index.js";

export type AssignableFeature =
  | "matchSummary"
  | "resumeGeneration"
  | "coverLetterGeneration";

export interface FeatureAssignment {
  provider: LlmProvider | "None";
  keyIndex: number;
}

export interface Settings {
  apiKeys: Partial<Record<LlmProvider, string[]>>;
  featureAssignments: Record<AssignableFeature, FeatureAssignment>;
}
