export type LlmProvider =
  | "Gemini"
  | "Cohere"
  | "OpenAI"
  | "Groq"
  | "Anthropic"
  | "Cerebras";

export type GeminiModels =
  | "gemini-2.0-flash" // Deprecated
  | "gemini-2.0-flash-lite" // Deprecated
  | "gemini-2.5-flash"
  | "gemini-2.5-flash-lite"
  | "gemini-2.5-pro"
  | "gemini-3.0-flash-preview"
  | "gemini-3.0-pro-preview";

export type CohereModels =
| "command-a-03-2025"
| "command-r7b-12-2024"
| "command-a-reasoning-08-2025"
| "command-r-08-2024"
| "command-r-plus-08-2024";

export type OpenAIModels =
| "gpt-5.2-2025-12-11"
| "gpt-5-mini-2025-08-07"
| "gpt-5-nano-2025-08-07"
| "gpt-5.2-pro-2025-12-11"
| "gpt-5-2025-08-07"
| "gpt-4.1-2025-04-14";

export type GroqModels =
| "llama-3.1-8b-instant"
| "llama-3.3-70b-versatile"
| "meta-llama/llama-guard-4-12b"
| "openai/gpt-oss-120b"
| "openai/gpt-oss-20b";

export type AnthropicModels =
| "claude-sonnet-4-5"
| "claude-haiku-4-5"
| "claude-opus-4-5";

export type CerebrasModels =
| "llama3.1-8b"
| "llama-3.3-70b"
| "gpt-oss-120b"
| "qwen-3-32b";