import * as coverLetterService from "../../src/services/cover_letter.js";
import * as resumeService from "../../src/services/resume.js";

import { CompilationRequestSchema } from "../../src/events/index.js";
import { generateIfNeeded } from "../../src/services/dispatcherService.js";

jest.mock("@documents/services/compileResume.js", () => ({
  compileResume: jest.fn().mockResolvedValue({
    success: true,
    downloadUrl: "/tmp/resume.pdf",
    changesUrl: "/tmp/resume.json",
  }),
}));

jest.mock("@documents/services/compileCoverLetter.js", () => ({
  compileCoverLetter: jest.fn().mockResolvedValue({
    success: true,
    downloadUrl: "/tmp/coverLetter.pdf",
    changesUrl: "/tmp/coverLetter.json",
  }),
}));

describe("generateIfNeeded service", () => {
  let validResumeRequest: any;
  let validCoverLetterRequest: any;

  beforeEach(() => {
    validResumeRequest = {
      jobID: 1,
      userID: "user123",
      docType: "resume",
      userInfo: {},
      educationInfo: {},
      resume: {},
      coverLetter: null,
      companyName: "Test Co",
    };

    validCoverLetterRequest = {
      ...validResumeRequest,
      docType: "cover-letter",
    };
  });

  it("should generate a resume when docType is 'resume'", async () => {
    const result = await generateIfNeeded(validResumeRequest);
    expect(result.success).toBe(true);
    expect(result.downloadUrl).toBe("/tmp/resume.pdf");
    expect(resumeService.compileResume).toHaveBeenCalledWith(
      validResumeRequest
    );
  });

  it("should generate a cover letter when docType is 'cover-letter'", async () => {
    const result = await generateIfNeeded(validCoverLetterRequest);
    expect(result.success).toBe(true);
    expect(result.downloadUrl).toBe("/tmp/coverLetter.pdf");
    expect(coverLetterService.compileCoverLetter).toHaveBeenCalledWith(
      validCoverLetterRequest
    );
  });

  it("should throw if docType is invalid", async () => {
    const invalidRequest = { ...validResumeRequest, docType: "invalid" };
    await expect(generateIfNeeded(invalidRequest)).rejects.toThrow();
  });
});
