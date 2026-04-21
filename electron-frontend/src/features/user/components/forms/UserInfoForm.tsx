import { FormProps } from "../../types";
import React from "react";
import { TextArea } from "@/components/UI/TextArea";
import { TextInput } from "@/components/UI/TextInput";
import { createSectionChangeHandler } from "../../utils/formUtils";

/**
 * A form component for editing the user's personal and contact information.
 * @param {FormProps} props The props for the component.
 * @returns {React.FC<FormProps>}
 */
export const UserInfoForm: React.FC<FormProps> = ({ profile, setProfile }) => {
  const handleChange = createSectionChangeHandler("userInfo", setProfile);

  return (
    <div className="card">
      <h2>Personal Information</h2>
      <div className="form-grid">
        <TextInput name="first_name" value={profile.userInfo.first_name} onChange={handleChange} placeholder="First Name" />
        <TextInput name="last_name" value={profile.userInfo.last_name} onChange={handleChange} placeholder="Last Name" />
        <TextInput name="email" value={profile.userInfo.email} onChange={handleChange} placeholder="Email" />
        <TextInput name="mobile" value={profile.userInfo.mobile} onChange={handleChange} placeholder="Mobile Phone" />
        <TextInput name="github" value={profile.userInfo.github} onChange={handleChange} placeholder="GitHub URL" />
        <TextInput name="linkedin" value={profile.userInfo.linkedin} onChange={handleChange} placeholder="LinkedIn URL" />
        <TextInput name="current_location" value={profile.userInfo.current_location} onChange={handleChange} placeholder="Current Location (e.g., City, State)" />
      </div>
      <TextArea name="summary" value={profile.userInfo.summary} onChange={handleChange} placeholder="Professional Summary..." rows={4} />
    </div>
  );
};