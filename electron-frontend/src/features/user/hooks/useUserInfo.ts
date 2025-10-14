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

const cleanLoadedExperiences = (experiences: Experience[] | undefined): Experience[] => {
  if (!experiences) return initialProfileState.resume.experiences || [];

  return experiences.map((exp) => {
    const cleanedExp = { ...exp };
    if (cleanedExp.bulletPoints) {
      cleanedExp.bulletPoints = cleanLoadedBulletPoints(
        cleanedExp.bulletPoints
      );
    }
    return cleanedExp;
  }) as Experience[];
};

const cleanLoadedProjects = (projects: Project[] | undefined): Project[] => {
  if (!projects) return initialProfileState.resume.projects || [];

  return projects.map((proj) => {
    const cleanedProj = { ...proj };
    if (cleanedProj.bulletPoints) {
      cleanedProj.bulletPoints = cleanLoadedBulletPoints(
        cleanedProj.bulletPoints
      );
    }
    return cleanedProj;
  }) as Project[];
};

const wasDataCleaned = (originalArray: Experience[] | Project[] | undefined, cleanedArray: Experience[] | Project[]): boolean => {
    if (!originalArray) return false;

    const countBulletPoints = (items: Experience[] | Project[]) => {
        return items.reduce((acc, item) => acc + (item.bulletPoints?.length || 0), 0);
    };

    const originalBulletPointCount = countBulletPoints(originalArray);
    const cleanedBulletPointCount = countBulletPoints(cleanedArray);

    return originalBulletPointCount !== cleanedBulletPointCount;
}

/**
 * A hook that provides the user profile and functions to save and load it.
 *
 * It returns an object with the following properties:
 * - `userProfile`: The current user profile.
 * - `setUserProfile`: A function to save the user profile.
 * - `loading`: A boolean indicating whether the user profile is currently being loaded or saved.
 * - `error`: A string indicating any error that occurred while loading or saving the user profile.
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
          throw new Error(
            profileResult.error || samplesResult.error || "Failed to load data."
          );
        }

        const loadedProfile = profileResult.data;
        const loadedSamples = samplesResult.data || [];

        if (!loadedProfile) {
          setUserProfile({
            ...initialProfileState,
            writingSamples: loadedSamples,
          });
        } else {
         const cleanedExperiences = cleanLoadedExperiences(loadedProfile.resume?.experiences);
          const cleanedProjects= cleanLoadedProjects(loadedProfile.resume?.projects);
          
          const originalExperiences = loadedProfile.resume?.experiences;
          const originalProjects = loadedProfile.resume?.projects;

          if (
            wasDataCleaned(originalExperiences, cleanedExperiences) ||
            wasDataCleaned(originalProjects, cleanedProjects)
          ) {
            cleanedData = true;
          }

          const mergedProfile: UserProfile = {
            userInfo: {
              ...initialProfileState.userInfo,
              ...(loadedProfile.userInfo || {}),
            },
            education: {
              ...initialProfileState.education,
              ...(loadedProfile.education || {}),
            },
            resume: {
              ...initialProfileState.resume,
              experiences: cleanedExperiences,
              projects: cleanedProjects,
              skills:
                loadedProfile.resume?.skills ||
                initialProfileState.resume.skills,
            },
            coverLetter: {
              ...initialProfileState.coverLetter,
              ...(loadedProfile.coverLetter || {}),
            },
            aboutMe: {
              ...initialProfileState.aboutMe,
              ...(loadedProfile.aboutMe || {}),
            },
            writingSamples: loadedSamples,
          };

          setUserProfile(mergedProfile);

          if (cleanedData) {
            console.warn(
              "Corrupted profile data found and cleaned during load. Resaving clean file."
            );
            await window.appAPI.user.saveUserInfo(mergedProfile);
          }
        }
      } catch (err: any) {
        setError(err.message ?? "An unknown error occurred.");
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, []);

  const saveUserProfile = useCallback(async (newProfile: UserProfile) => {
    setLoading(true);
    setError(null);

    try {
      const { writingSamples, ...profileToSave } = newProfile;
      const result = await window.appAPI.user.saveUserInfo(profileToSave);

      if (result.success) {
        setUserProfile(newProfile);
      } else {
        setError(result.error ?? "Failed to save profile data on disk.");
      }
    } catch (err: any) {
      setError(
        err.message ?? "An unexpected error occurred during profile save."
      );
      console.error("Error saving user profile:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  return { userProfile, setUserProfile, loading, error, saveUserProfile };
};