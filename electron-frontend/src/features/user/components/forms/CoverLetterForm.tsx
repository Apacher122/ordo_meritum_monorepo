import React from "react";
import { UserProfile } from "../../types";

interface FormProps {
  profile: UserProfile;
  setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>;
}

export const CoverLetterForm: React.FC<FormProps> = ({ profile, setProfile }) => {
    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setProfile(prev => prev && { ...prev, coverLetter: { ...prev.coverLetter, [name]: value } });
    };

    return (
        <div className="card form-vertical">
            <h2>Cover Letter Sections</h2>
            <label>About</label>
            <textarea name="about" value={profile.coverLetter.about || ""} onChange={handleChange} placeholder="Write the 'About' section of your cover letter..." className="textarea" rows={6} />
            <label>Experience</label>
            <textarea name="experience" value={profile.coverLetter.experience || ""} onChange={handleChange} placeholder="Write the 'Experience' section..." className="textarea" rows={6} />
            <label>What I Bring</label>
            <textarea name="whatIBring" value={profile.coverLetter.whatIBring || ""} onChange={handleChange} placeholder="Write the 'What I Bring' section..." className="textarea" rows={6} />
        </div>
    );
};