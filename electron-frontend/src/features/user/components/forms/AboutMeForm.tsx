import React from "react";
import { UserProfile } from "../../types";

interface FormProps {
  profile: UserProfile;
  setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>;
}

export const AboutMeForm: React.FC<FormProps> = ({ profile, setProfile }) => {
    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const { value } = e.target;
        setProfile(prev => prev && { ...prev, aboutMe: { ...prev.aboutMe, essay: value } });
    };

    return (
        <div className="card form-vertical">
            <h2>About Me</h2>
            <textarea name="essay" value={profile.aboutMe.essay || ""} onChange={handleChange} placeholder="Write your 'About Me' essay here..." className="textarea" rows={15} />
        </div>
    );
};