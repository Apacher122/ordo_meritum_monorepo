import { FormProps } from "../../types";
import React from "react";
import { TextArea } from "@/components/UI/TextArea";
import { createSectionChangeHandler } from "../../utils/formUtils";

/**
* A form component for editing the "About Me" essay section of the user profile.
* @param {FormProps} props The props for the component.
* @returns {React.FC<FormProps>}
*/
export const AboutMeForm: React.FC<FormProps> = ({ profile, setProfile }) => {
  const handleChange = createSectionChangeHandler("aboutMe", setProfile);
  
  return (
    <div className="card form-vertical">
    <h2>About Me</h2>
    <TextArea name="essay" value={profile.aboutMe.essay} onChange={handleChange} placeholder="Write your 'About Me' essay here..." rows={15} />
    </div>
  );
};