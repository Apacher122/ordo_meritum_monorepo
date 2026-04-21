import { FormProps } from "../../types";
import React from "react";
import { TextArea } from "@/components/UI/TextArea";
import { createSectionChangeHandler } from "../../utils/formUtils";

/**
* A form component for editing the distinct sections of a generic cover letter.
* @param {FormProps} props The props for the component.
* @returns {React.FC<FormProps>}
*/
export const CoverLetterForm: React.FC<FormProps> = ({ profile, setProfile }) => {
  const handleChange = createSectionChangeHandler("coverLetter", setProfile);
  
  return (
    <div className="card form-vertical">
    <h2>Cover Letter Sections</h2>
    <label htmlFor="coverLetter-about">About</label>
    <TextArea id="coverLetter-about" name="about" value={profile.coverLetter.about} onChange={handleChange} placeholder="Write the 'About' section of your cover letter..." rows={6} />
    <label htmlFor="coverLetter-experience">Experience</label>
    <TextArea id="coverLetter-experience" name="experience" value={profile.coverLetter.experience} onChange={handleChange} placeholder="Write the 'Experience' section..." rows={6} />
    <label htmlFor="coverLetter-whatIBring">What I Bring</label>
    <TextArea id="coverLetter-whatIBring" name="whatIBring" value={profile.coverLetter.whatIBring} onChange={handleChange} placeholder="Write the 'What I Bring' section..." rows={6} />
    </div>
  );
};