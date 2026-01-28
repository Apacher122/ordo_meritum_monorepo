export class ResumeBuilderError extends Error {
  timestamp: string;
  details: Record<string, any>;
  file?: string;
  line?: number;

/**
 * Constructor for ResumeBuilderError.
 *
 * @param {string} message - The error message.
 * @param {Record<string, any>} details - Optional details about the error.
 *
 * This constructor sets the error message, timestamp, and details.
 * It also captures the current stack trace if available.
 * If details.captureStackTrace is true, it sets the file and line number where the error occurred.
 */
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
