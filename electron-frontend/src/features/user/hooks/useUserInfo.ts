import { useCallback, useEffect, useState } from "react";

import { UserProfile } from "../types";
import { app } from '../../../config/firebase';

const initialProfileState: UserProfile = {
  userInfo: {},
  education: {},
  resume: { skills: [{skill: ""}], experiences: [], projects: [] },
  coverLetter: {},
  aboutMe: {},
  writingSamples: [],
};


/**
 * A hook that provides the user profile and functions to save and load it.
 *
 * It returns an object with the following properties:
 *   - `userProfile`: The current user profile.
 *   - `setUserProfile`: A function to save the user profile.
 *   - `loading`: A boolean indicating whether the user profile is currently being loaded or saved.
 *   - `error`: A string indicating any error that occurred while loading or saving the user profile.
 */
export const useUserInfo = () => {
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      setError(null);
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
              ...(loadedProfile.resume || {}),
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
    const { writingSamples, ...profileToSave } = newProfile;
    const result = await window.appAPI.user.saveUserInfo(profileToSave);
    if (result.success) {
      setUserProfile(newProfile);
    } else {
      setError(result.error ?? null);
    }
    setLoading(false);
  }, []);

  return { userProfile, setUserProfile, loading, error, saveUserProfile };
};
