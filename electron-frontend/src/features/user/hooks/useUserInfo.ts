import { Experience, Project, UserProfile } from "../types";
import { useCallback, useEffect, useState } from "react";

const initialProfileState: UserProfile = {
  userInfo: {},
  education: {},
  resume: { skills: [{ skill: "" }], experiences: [], projects: [] },
  coverLetter: {},
  aboutMe: {},
  writingSamples: [],
};

const cleanLoadedBulletPoints = (bulletPoints: any[] | undefined) => {
  if (!Array.isArray(bulletPoints)) return [];

  return bulletPoints
    .map((bp) => {
      if (typeof bp?.text !== "string") {
        return null;
      }

      const cleanedBp: { text: string; id?: string } = {
        text: bp.text,
      };

      if (typeof bp.id === "string") {
        cleanedBp.id = bp.id;
      }

      return cleanedBp;
    })
    .filter((bp) => bp !== null) as { text: string; id?: string }[];
};

const cleanResumeItems = <T extends Experience | Project>(
  items: T[] | undefined,
  fallback: T[]
): T[] => {
  if (!items) return fallback;
  return items.map((item) => ({
    ...item,
    bulletPoints: item.bulletPoints ? cleanLoadedBulletPoints(item.bulletPoints) : []
  })) as T[];
};

const wasDataCleaned = (
  originalArray: Experience[] | Project[] | undefined,
  cleanedArray: Experience[] | Project[]
): boolean => {
  if (!originalArray) return false;
  const countBulletPoints = (items: Experience[] | Project[]) => 
    items.reduce((acc, item) => acc + (item.bulletPoints?.length || 0), 0);
    
  return countBulletPoints(originalArray) !== countBulletPoints(cleanedArray);
};

/**
 * Custom hook to manage the user's profile data.
 * Handles fetching the profile from a persistent source, merging it with initial state,
 * cleaning potentially corrupted data, and providing a function to save updates.
 * @returns {UseUserInfoReturn} An object containing the user profile state and management functions.
 */
export const useUserInfo = () => {
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      setError(null);
      let cleanedData = false;

      try {
        const [profileResult, samplesResult] = await Promise.all([
          window.appAPI.user.loadUserInfo(),
          window.appAPI.writingSamples.load(),
        ]);

        if (!profileResult.success || !samplesResult.success) {
          throw new Error(profileResult.error || samplesResult.error || "Failed to load data.");
        }

        const loadedProfile = profileResult.data;
        const loadedSamples = samplesResult.data || [];

        if (loadedProfile) {
          const cleanedExperiences = cleanResumeItems(
            loadedProfile.resume?.experiences, 
            initialProfileState.resume.experiences || []
          );
          const cleanedProjects = cleanResumeItems(
            loadedProfile.resume?.projects,
            initialProfileState.resume.projects || []
          );

          if (
            wasDataCleaned(loadedProfile.resume?.experiences, cleanedExperiences) ||
            wasDataCleaned(loadedProfile.resume?.projects, cleanedProjects)
          ) {
            cleanedData = true;
          }

          const mergedProfile: UserProfile = {
            userInfo: { ...initialProfileState.userInfo, ...loadedProfile.userInfo },
            education: { ...initialProfileState.education, ...loadedProfile.education },
            resume: {
              ...initialProfileState.resume,
              experiences: cleanedExperiences,
              projects: cleanedProjects,
              skills: loadedProfile.resume?.skills || initialProfileState.resume.skills,
            },
            coverLetter: { ...initialProfileState.coverLetter, ...loadedProfile.coverLetter },
            aboutMe: { ...initialProfileState.aboutMe, ...loadedProfile.aboutMe },
            writingSamples: loadedSamples,
          };

          setUserProfile(mergedProfile);

          if (cleanedData) {
            console.warn("Corrupted profile data found and cleaned during load. Resaving clean file.");
            await window.appAPI.user.saveUserInfo(mergedProfile);
          }
        } else {
          setUserProfile({ ...initialProfileState, writingSamples: loadedSamples });
        }
      } catch (err: any) {
        setError(err.message ?? "An unknown error occurred.");
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, []);

  const saveUserProfile = useCallback(
    async (newProfile: UserProfile) => {
      const previousProfile = userProfile;
      setUserProfile(newProfile);
      setError(null);

      try {
        const { writingSamples, ...profileToSave } = newProfile;
        const result = await window.appAPI.user.saveUserInfo(profileToSave);

        if (!result.success) {
          setUserProfile(previousProfile);
          setError(result.error ?? "Failed to save profile data on disk.");
        }
      } catch (err: any) {
        setUserProfile(previousProfile);
        setError(err.message ?? "An unexpected error occurred during profile save.");
      }
    },
    [userProfile]
  );

  return { userProfile, setUserProfile, loading, error, saveUserProfile };
};