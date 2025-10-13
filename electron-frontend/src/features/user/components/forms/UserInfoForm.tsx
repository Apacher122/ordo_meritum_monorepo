import React from "react";
import { UserProfile } from "../../types";

interface FormProps {
  profile: UserProfile;
  setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>;
}

export const UserInfoForm: React.FC<FormProps> = ({ profile, setProfile }) => {
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setProfile(
      (prev) =>
        prev && {
          ...prev,
          userInfo: { ...prev.userInfo, [name]: value },
        }
    );
  };

  return (
    <div className="card">
      <h2>Personal Information</h2>
      <div className="form-grid">
        <input name="first_name" value={profile.userInfo.first_name || ""} onChange={handleChange} placeholder="First Name" className="input" />
        <input name="last_name" value={profile.userInfo.last_name || ""} onChange={handleChange} placeholder="Last Name" className="input" />
        <input name="email" value={profile.userInfo.email || ""} onChange={handleChange} placeholder="Email" className="input" />
        <input name="mobile" value={profile.userInfo.mobile || ""} onChange={handleChange} placeholder="Mobile Phone" className="input" />
        <input name="github" value={profile.userInfo.github || ""} onChange={handleChange} placeholder="GitHub URL" className="input" />
        <input name="linkedin" value={profile.userInfo.linkedin || ""} onChange={handleChange} placeholder="LinkedIn URL" className="input" />
        <input name="current_location" value={profile.userInfo.current_location || ""} onChange={handleChange} placeholder="Current Location (e.g., City, State)" className="input" />
      </div>
      <textarea name="summary" value={profile.userInfo.summary || ""} onChange={handleChange} placeholder="Professional Summary..." className="textarea" rows={4} />
    </div>
  );
};