import {
  CompilationRequestSchema,
  CoverLetterPayloadSchema,
} from "../../src/events/index.js";

import { compileCoverLetter } from "../../src/services/cover_letter.js";

jest.mock("@documents/utils/documents", () => ({
  initializeDocumentWorkspace: jest.fn().mockResolvedValue({
    tempFolder: "/tmp",
    tempPdf: "/tmp/coverLetter.pdf",
    tempFolderCompiled: "/tmp/compiled",
    tempJson: "/tmp/data.json",
  }),
  createHeader: jest.fn(),
  saveJson: jest.fn().mockResolvedValue("/tmp/data.json"),
}));

jest.mock("@documents/services/export.js", () => ({
  exportLatex: jest.fn().mockResolvedValue("/tmp/coverLetter.pdf"),
}));

describe("compileCoverLetter service", () => {
  let validRequest: any;

  beforeEach(() => {
    validRequest = {
      jobID: 1,
      userID: "user123",
      companyName: "Test Co",
      docType: "cover-letter",
      userInfo: {
        first_name: "Jane",
        last_name: "Doe",
        current_location: "NY",
        email: "jane@test.com",
        github: "jane",
        linkedin: "jane",
        mobile: "1234567890",
        summary: "Experienced developer",
      },
      educationInfo: {
        coursework: "CS",
        degree: "BS",
        location: "NY",
        school: "Test University",
        start_end: "2010-2014",
      },
      resume: {},
      coverLetter: {
        companyName: "Test Co",
        companyProperName: "Test Company",
        jobTitle: "Software Engineer",
        body: {
          about: "I am a great candidate.",
          experience: "Worked at X and Y.",
          whatIBring: "Skills in TS and Go.",
          revisionSummary: "First draft",
        },
      },
    };

    validRequest = CompilationRequestSchema.parse(validRequest);
  });

  it("should successfully compile a valid cover letter request", async () => {
    const result = await compileCoverLetter(validRequest);
    expect(result.success).toBe(true);
    expect(result.downloadUrl).toBe("/tmp/coverLetter.pdf");
    expect(result.changesUrl).toBe("/tmp/data.json");
  });

  it("should throw Zod validation error for invalid email", () => {
    const invalidRequest = {
      ...validRequest,
      userInfo: { ...validRequest.userInfo, email: "not-an-email" },
    };
    expect(() => CompilationRequestSchema.parse(invalidRequest)).toThrow();
  });

  it("should handle errors gracefully", async () => {
    const { initializeDocumentWorkspace } = await import(
      "../../src/utils/documents/workspace.js"
    );
    (initializeDocumentWorkspace as jest.Mock).mockRejectedValueOnce(
      new Error("FS error")
    );

    const result = await compileCoverLetter(validRequest);
    expect(result.success).toBe(false);
    expect(result.error).toBe("FS error");
  });
});
