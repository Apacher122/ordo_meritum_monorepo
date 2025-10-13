import { CompilationRequestSchema } from "../../src/events";

const validRequest = {
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

export const parsedRequest = CompilationRequestSchema.parse(validRequest);
