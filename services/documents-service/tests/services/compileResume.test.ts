import { CompilationRequestSchema } from "../../src/events/index.js";
import { compileResume } from "../../src/services/resume.js";

jest.mock("@documents/utils/documents", () => ({
  initializeDocumentWorkspace: jest.fn().mockResolvedValue({
    tempFolder: "/tmp",
    tempPdf: "/tmp/test.pdf",
    tempFolderCompiled: "/tmp/compiled",
    tempJson: "/tmp/data.json",
  }),
  createHeader: jest.fn(),
  saveJson: jest.fn().mockResolvedValue("/tmp/data.json"),
}));

jest.mock("@documents/utils/latex", () => ({
  generateLatexSectionFile: jest.fn(),
}));

jest.mock("@documents/services/export", () => ({
  exportLatex: jest.fn().mockResolvedValue("/tmp/test.pdf"),
}));

describe("compileResume service", () => {
  let validRequest: any;

  beforeEach(() => {
    validRequest = {
      jobID: 1,
      userID: "user123",
      companyName: "Test Co",
      docType: "resume",
      userInfo: {
        first_name: "John",
        last_name: "Doe",
        current_location: "NY",
        email: "john@test.com",
        github: "john",
        linkedin: "john",
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
      resume: {
        summary: [
          {
            sentence: "Test summary",
            justification_for_change: "",
            is_new_suggestion: true,
          },
        ],
        skills: [],
        experiences: [],
        projects: [],
      },
      coverLetter: null,
    };

    validRequest = CompilationRequestSchema.parse(validRequest);
  });

  it("should successfully compile a valid resume request", async () => {
    const result = await compileResume(validRequest);
    expect(result.success).toBe(true);
    expect(result.downloadUrl).toBe("/tmp/test.pdf");
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
    const { initializeDocumentWorkspace } = await import("../../src/utils/documents/index.js");
    (initializeDocumentWorkspace as jest.Mock).mockRejectedValueOnce(new Error("FS error"));

    const result = await compileResume(validRequest);
    expect(result.success).toBe(false);
    expect(result.error).toBe("FS error");
  });
});
