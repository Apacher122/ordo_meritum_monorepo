export class ResumeBuilderError extends Error {
  timestamp: string;
  details: Record<string, any>;
  file?: string;
  line?: number;

  constructor(message: string, details: Record<string, any> = {}) {
    super(message);
    this.name = this.constructor.name;
    this.timestamp = new Date().toISOString();
    this.details = details;

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }

    if (details.captureStackTrace) {
      const stack = new Error().stack?.split('\n')[2];
      if (stack) {
        const match = stack.match(/\((.*?):(\d+):\d+\)/);
        if (match) {
          this.file = match[1];
          this.line = parseInt(match[2], 10);
        }
      }
    }
  }
}

export class LaTeXFileAccessError extends ResumeBuilderError {
  constructor(
    message = "Error accessing LaTeX file",
    details: Record<string, any> = {},
  ) {
    super(message, details);
  }
}

export class ResumeSectionNotFoundError extends ResumeBuilderError {
  constructor(message = "Resume section not found", details: Record<string, any> = {}) {
    super(message, details);
  }
}

export class RateLimitError extends ResumeBuilderError {
  constructor(message = "Rate limit exceeded", details: Record<string, any> = {}) {
    super(message,details);
  }
}
