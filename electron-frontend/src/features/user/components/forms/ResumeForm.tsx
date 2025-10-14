import { Experience, Project, UserProfile } from "../../types";

import { BulletPointInput } from "../BulletPointInput";
import React from "react";

interface FormProps {
  profile: UserProfile;
  setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>;
}

export const ResumeForm: React.FC<FormProps> = ({ profile, setProfile }) => {

  const handleSkillsChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const { value } = e.target;
      setProfile(prev => prev && { ...prev, resume: { ...prev.resume, skills: value as any }});
  };

  const handleExperienceChange = (index: number, field: keyof Experience, value: any) => {
    setProfile(prev => {
        if (!prev) return null;
        const newExperiences = [...(prev.resume.experiences || [])];
        const updatedValue = field === 'bulletPoints' ? value as { text: string }[] : value;     
        newExperiences[index] = { ...newExperiences[index], [field]: updatedValue };
        return { ...prev, resume: { ...prev.resume, experiences: newExperiences } };
    });
  };

  const addExperience = () => {
    setProfile(prev => {
        if (!prev) return null;
        const newExperience: Experience = { id: Date.now().toString(), bulletPoints: [{ text: "" }] };
        const experiences = [...(prev.resume.experiences || []), newExperience];
        return { ...prev, resume: { ...prev.resume, experiences } };
    });
  };

  
  const handleProjectChange = (index: number, field: keyof Project, value: any) => {
    setProfile(prev => {
        if (!prev) return null;
        const newProjects = [...(prev.resume.projects || [])];
        const updatedValue = field === 'bulletPoints' ? value as { text: string }[] : value;

        newProjects[index] = { ...newProjects[index], [field]: updatedValue };
        return { ...prev, resume: { ...prev.resume, projects: newProjects } };
    });
  };

  const addProject = () => {
    setProfile(prev => {
        if (!prev) return null;
        const newProject: Project = { id: Date.now().toString(), bulletPoints: [{ text: "" }] };
        const projects = [...(prev.resume.projects || []), newProject];
        return { ...prev, resume: { ...prev.resume, projects } };
    });
  };

  return (
    <div className="resume-form">
      <div className="card">
        <h2>Skills</h2>
        <textarea name="skills" value={profile.resume.skills as unknown as string || ""} onChange={handleSkillsChange} placeholder="List your skills, separated by commas..." className="textarea" rows={4} />
      </div>

      <div className="card">
        <h2>Work Experience</h2>
        {(profile.resume.experiences || []).map((exp, index) => (
          <div key={exp.id} className="experience-entry">
            <input value={exp.company || ""} onChange={(e) => handleExperienceChange(index, 'company', e.target.value)} placeholder="Company" className="input" />
            <input value={exp.jobTitle || ""} onChange={(e) => handleExperienceChange(index, 'jobTitle', e.target.value)} placeholder="Job Title" className="input" />
            <input value={exp.years || ""} onChange={(e) => handleExperienceChange(index, 'years', e.target.value)} placeholder="Years of Employment" className="input" />
            <BulletPointInput bullets={exp.bulletPoints || []} onChange={(bullets) => handleExperienceChange(index, 'bulletPoints', bullets)} />
          </div>
        ))}
        <button type="button" onClick={addExperience} className="button">Add Experience</button>
      </div>

       <div className="card">
        <h2>Projects</h2>
        {(profile.resume.projects || []).map((proj, index) => (
          <div key={proj.id} className="experience-entry">
            <input value={proj.name || ""} onChange={(e) => handleProjectChange(index, 'name', e.target.value)} placeholder="Project Name" className="input" />
            <input value={proj.description || ""} onChange={(e) => handleProjectChange(index, 'description', e.target.value)} placeholder="Project Description" className="input" />
            <input value={proj.years || ""} onChange={(e) => handleProjectChange(index, 'years', e.target.value)} placeholder="Years" className="input" />
            <BulletPointInput bullets={proj.bulletPoints || []} onChange={(bullets) => handleProjectChange(index, 'bulletPoints', bullets)} />
          </div>
        ))}
        <button type="button" onClick={addProject} className="button">Add Project</button>
      </div>
    </div>
  );
};