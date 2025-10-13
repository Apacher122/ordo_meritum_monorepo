import * as user from "@/features/user/types/types";

export interface DocumentRequestBody {
  userInfo: user.UserInfo;
  education: user.EducationInfo;
  resume: user.ResumeData;
  coverLetter: user.CoverLetterData;
  aboutMe: user.AboutMeData;
  writingSamples?: user.WritingSample[];
}