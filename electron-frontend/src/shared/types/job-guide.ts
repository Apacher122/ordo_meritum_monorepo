export interface ScoreMetric {
  score_title: string;
  raw_score: number;
  weighted_score: number;
  score_weight: number;
  score_reason: string;
  isCompatible: boolean;
  strength: string;
  weaknesses: string;
}

export interface SummaryInfo {
  summary_text: string;
  summary_temperature: number;
}

export interface OverallSummary {
  overall_match_score: number;
  summary: SummaryInfo[];
  suggestions: string[];
}

export interface MatchSummary {
  id: number;
  job_posting_id: number;
  should_apply: boolean;
  should_apply_reasoning: string;
  metrics: ScoreMetric[];
  overall_match_summary: OverallSummary;
  created_at: string;
  updated_at: string;
}

export interface GuidingQuestion {
  question: string;
  answer: string;
  suggestions: string[];
}

export interface RawMatchSummary {
  id: number;
  job_posting_id: number;
  should_apply: string;
  should_apply_reasoning: string;
  metrics_json: string;
  overall_summary_json: string;
  projects_section_missing_entries: boolean;
  created_at: string;
  updated_at: string;
}

export interface RawMatchSummaryResponse {
  success: boolean;
  matchSummary: RawMatchSummary;
}

export interface GuidingQuestionsResponse {
  success: boolean;
  guidingQuestions: GuidingQuestion[];
}
