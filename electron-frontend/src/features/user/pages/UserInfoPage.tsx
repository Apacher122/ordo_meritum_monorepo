import "@/assets/styles/pages/UserInfoPage.css";

import { Experience, Project, UserProfile } from "../types";
import { ProfileTab, TabHeaderControls } from "../components/TabHeaderControls";
import React, { useEffect, useMemo, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { AboutMeForm } from "../components/forms/AboutMeForm";
import { CoverLetterForm } from "../components/forms/CoverLetterForm";
import { EducationForm } from "../components/forms/EducationForm";
import { ResumeForm } from "../components/forms/ResumeForm";
import { UploadSamplesButton } from "../components/UploadSamplesButton";
import { UserInfoForm } from "../components/forms/UserInfoForm";
import { useUserInfo } from "../hooks/useUserInfo";

const isBulletPointValid = (bp: { text: string } | null | undefined): boolean => {
    if (!bp || typeof bp.text !== 'string') {
        return false;
    }
    return bp.text.trim().length > 0; 
};

const hasValidBulletPoints = (item: Experience | Project): boolean => {
    return (item.bulletPoints?.length ?? 0) > 0;
};

const cleanBulletPoints = <T extends Experience | Project>(items: T[] | undefined): T[] => {
    if (!items) return [];
    
    const itemsWithCleanedBullets = items.map(item => ({
        ...item,
        bulletPoints: item.bulletPoints?.filter(isBulletPointValid)
    }));
    
    return itemsWithCleanedBullets.filter(hasValidBulletPoints) as T[];
};

export const UserInfoPage: React.FC = () => {
  const { userProfile, loading, error, saveUserProfile } = useUserInfo();
  const [activeTab, setActiveTab] = useState<ProfileTab>("User Info");
  const [formState, setFormState] = useState<UserProfile | null>(null);
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();
  
  const headerControls = useMemo(() => {
    const tabs = (
      <TabHeaderControls activeTab={activeTab} onTabChange={setActiveTab} />
    );
    
    if (activeTab === "User Info") {
      return (
        <div className="header-controls-wrapper">
        {tabs}
        <UploadSamplesButton />
        </div>
      );
    }
    return tabs;
  }, [activeTab]);
  
  useEffect(() => {
    setHeaderTitle("Your Information");
    setHeaderSubtitle("Manage your personal details, resume, and more");
    setHeaderControls(headerControls);
    
    return () => {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
      setHeaderControls(null);
    };
  }, [setHeaderTitle, setHeaderSubtitle, setHeaderControls, headerControls]);
  
  
  useEffect(() => {
    if (userProfile) {
      const skillsArray = userProfile.resume.skills as { skill: string }[] | undefined;
      
      const skillsString = (skillsArray && Array.isArray(skillsArray))
      ? skillsArray
      .map(s => (s && typeof s.skill === 'string') ? s.skill : '') 
      .filter(Boolean)
      .join(', ')
      : '';
      
      setFormState({
        ...userProfile,
        resume: {
          ...userProfile.resume,
          skills: skillsString as any, 
        },
      });
    }
  }, [userProfile]);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (formState) {
    const skillsString = (formState.resume.skills as unknown as string) || "";
    const skillsArray = skillsString
      .split(',')
      .map(skill => skill.trim())
      .filter(skill => skill)
      .map(skillName => ({ skill: skillName })); 

    
    const cleanedExperiences = cleanBulletPoints(formState.resume.experiences);
    const cleanedProjects = cleanBulletPoints(formState.resume.projects);
      let payload: UserProfile;
      payload = {
        ...formState,
        resume: {
          ...formState.resume,
          skills: skillsArray,
          experiences: cleanedExperiences,
          projects: cleanedProjects,
        },
      };
      
      await saveUserProfile(payload);
      alert("Profile Saved!"); 
    }
  };
  
  const renderActiveTab = () => {
    if (!formState) return null;
    
    switch (activeTab) {
      case "User Info":
      return <UserInfoForm profile={formState} setProfile={setFormState} />;
      case "Education":
      return <EducationForm profile={formState} setProfile={setFormState} />;
      case "Resume":
      return <ResumeForm profile={formState} setProfile={setFormState} />;
      case "Cover Letter":
      return (
        <CoverLetterForm profile={formState} setProfile={setFormState} />
      );
      case "About Me":
      return <AboutMeForm profile={formState} setProfile={setFormState} />;
      default:
      return <div>Please select a section to edit.</div>;
    }
  };
  
  if (loading || !formState) {
    return <div>Loading user information...</div>;
  }
  
  return (
    <div className="user-info-page">
    <form onSubmit={handleSubmit} className="info-form">
    {error && <div className="error-message">{error}</div>}
    
    <div className="form-content">{renderActiveTab()}</div>
    
    <button type="submit" className="submit-button" disabled={loading}>
    {loading ? "Saving..." : "Save All Profile Information"}
    </button>
    </form>
    </div>
  );
};