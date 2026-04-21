import { BulletPointInput } from "../BulletPointInput";
import { FormProps } from "../../types";
import React from "react";
import { TextArea } from "@/components/UI/TextArea";
import { TextInput } from "@/components/UI/TextInput";

/**
 * A form for editing the core components of a resume, including skills,
 * work experiences, and projects.
 * @param {FormProps} props The props for the component.
 * @returns {React.FC<FormProps>}
 */
export const ResumeForm: React.FC<FormProps> = ({ profile, setProfile }) => {

  const handleSkillsChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const { value } = e.target;
      setProfile(prev => prev && { ...prev, resume: { ...prev.resume, skills: value as any }});
  };
  const handleItemChange = (index: number, arrayKey: 'experiences' | 'projects', field: string, value: any) => {
    setProfile(prev => {
        if (!prev) return null;
        const newItems = [...(prev.resume[arrayKey] || [])];
        const updatedValue = field === 'bulletPoints' ? value as { text: string }[] : value;     
        newItems[index] = { ...newItems[index], [field]: updatedValue };
        return { ...prev, resume: { ...prev.resume, [arrayKey]: newItems } };
    });
  };

  const addItem = (arrayKey: 'experiences' | 'projects') => {
    setProfile(prev => {
        if (!prev) return null;
        const newItem = { id: Date.now().toString(), bulletPoints: [{ text: "" }] } as any; 
        const items = [...(prev.resume[arrayKey] || []), newItem];
        return { ...prev, resume: { ...prev.resume, [arrayKey]: items } };
    });
  };

  return (
    <div className="resume-form">
      <div className="card">
        <h2>Skills</h2>
        <TextArea name="skills" value={profile.resume.skills as unknown as string} onChange={handleSkillsChange} placeholder="List your skills, separated by commas..." rows={4} />
      </div>

      <div className="card">
        <h2>Work Experience</h2>
        {(profile.resume.experiences || []).map((exp, index) => (
          <div key={exp.id} className="experience-entry">
            <TextInput value={exp.company} onChange={(e) => handleItemChange(index, 'experiences', 'company', e.target.value)} placeholder="Company" />
            <TextInput value={exp.jobTitle} onChange={(e) => handleItemChange(index, 'experiences', 'jobTitle', e.target.value)} placeholder="Job Title" />
            <TextInput value={exp.years} onChange={(e) => handleItemChange(index, 'experiences', 'years', e.target.value)} placeholder="Years of Employment" />
            <BulletPointInput bullets={exp.bulletPoints || []} onChange={(bullets) => handleItemChange(index, 'experiences', 'bulletPoints', bullets)} />
          </div>
        ))}
        <button type="button" onClick={() => addItem('experiences')} className="button">Add Experience</button>
      </div>

       <div className="card">
        <h2>Projects</h2>
        {(profile.resume.projects || []).map((proj, index) => (
          <div key={proj.id} className="experience-entry">
            <TextInput value={proj.name} onChange={(e) => handleItemChange(index, 'projects', 'name', e.target.value)} placeholder="Project Name" />
            <TextInput value={proj.description} onChange={(e) => handleItemChange(index, 'projects', 'description', e.target.value)} placeholder="Project Description" />
            <TextInput value={proj.years} onChange={(e) => handleItemChange(index, 'projects', 'years', e.target.value)} placeholder="Years" />
            <BulletPointInput bullets={proj.bulletPoints || []} onChange={(bullets) => handleItemChange(index, 'projects', 'bulletPoints', bullets)} />
          </div>
        ))}
        <button type="button" onClick={() => addItem('projects')} className="button">Add Project</button>
      </div>
    </div>
  );
};