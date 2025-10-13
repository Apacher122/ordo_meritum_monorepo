import React from "react";
import { UserProfile } from "../../types";

interface FormProps {
  profile: UserProfile;
  setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>;
}

export const EducationForm: React.FC<FormProps> = ({ profile, setProfile }) => {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setProfile(prev => prev && { ...prev, education: { ...prev.education, [name]: value } });
    };

    return (
        <div className="card">
            <h2>Education</h2>
            <div className="form-grid">
                <input name="school" value={profile.education.school || ""} onChange={handleChange} placeholder="School / University" className="input" />
                <input name="degree" value={profile.education.degree || ""} onChange={handleChange} placeholder="Degree (e.g., B.S. in Computer Science)" className="input" />
                <input name="start_end" value={profile.education.start_end || ""} onChange={handleChange} placeholder="Start - End Dates (e.g., Aug 2020 - May 2024)" className="input" />
                <input name="location" value={profile.education.location || ""} onChange={handleChange} placeholder="School Location (e.g., City, State)" className="input" />
            </div>
            <textarea name="coursework" value={profile.education.coursework || ""} onChange={handleChange} placeholder="Relevant Coursework..." className="textarea" rows={4} />
        </div>
    );
};