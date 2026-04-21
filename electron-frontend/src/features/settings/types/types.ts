import { LlmProvider } from "@/shared/types/index.js";

export type AssignableFeature =
  | "matchSummary"
  | "resumeGeneration"
  | "coverLetterGeneration";

export interface FeatureAssignment {
  provider: LlmProvider | "None";
  keyIndex: number;
}

export interface RateLimitConfig {
  callsPerDay: number;
  callsPerMinute: number;
  currentDayCount: number;
  currentMinuteCount: number;
  lastDayReset: number;
  lastMinuteReset: number;
}

export interface ApiKeyConfig {
  key: string;
  rateLimit: RateLimitConfig;
}

export interface Settings {
  apiKeys: Partial<Record<LlmProvider, ApiKeyConfig[]>>;
  featureAssignments: Record<AssignableFeature, FeatureAssignment>;
}
