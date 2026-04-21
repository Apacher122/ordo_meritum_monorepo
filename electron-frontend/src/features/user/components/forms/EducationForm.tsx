import { FormProps } from "../../types";
import React from "react";
import { TextArea } from "@/components/UI/TextArea";
import { TextInput } from "@/components/UI/TextInput";
import { createSectionChangeHandler } from "../../utils/formUtils";

/**
* A form component for editing the education section of the user profile.
* @param {FormProps} props The props for the component.
* @returns {React.FC<FormProps>}
*/
export const EducationForm: React.FC<FormProps> = ({ profile, setProfile }) => {
  const handleChange = createSectionChangeHandler("education", setProfile);
  
  return (
    <div className="card">
    <h2>Education</h2>
    <div className="form-grid">
    <TextInput name="school" value={profile.education.school} onChange={handleChange} placeholder="School / University" />
    <TextInput name="degree" value={profile.education.degree} onChange={handleChange} placeholder="Degree (e.g., B.S. in Computer Science)" />
    <TextInput name="start_end" value={profile.education.start_end} onChange={handleChange} placeholder="Start - End Dates (e.g., Aug 2020 - May 2024)" />
    <TextInput name="location" value={profile.education.location} onChange={handleChange} placeholder="School Location (e.g., City, State)" />
    </div>
    <TextArea name="coursework" value={profile.education.coursework} onChange={handleChange} placeholder="Relevant Coursework..." rows={4} />
    </div>
  );
};