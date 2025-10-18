import { useCallback, useEffect, useState } from "react";

import { Settings } from "../types/types";

const initialSettingsState: Settings = {
  apiKeys: {},
  featureAssignments: {
    matchSummary: "Gemini",
    resumeGeneration: "Gemini",
    coverLetterGeneration: "Gemini",
  },
};

export const useSettings = () => {
  const [settings, setSettings] = useState<Settings | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      setError(null);
      const result = await window.appAPI.user.loadSettings();
      if (result.success) {
        const loadedData = result.data;

        if (loadedData) {
          const mergedSettings: Settings = {
            apiKeys: {
              ...initialSettingsState.apiKeys,
              ...loadedData.apiKeys,
            },
            featureAssignments: {
              ...initialSettingsState.featureAssignments,
              ...loadedData.featureAssignments,
            },
          };
          setSettings(mergedSettings);
        } else {
          setSettings(initialSettingsState);
        }
      } else {
        setError(result.error ?? null);
      }
      setLoading(false);
    };
    loadData();
  }, []);

  const saveSettings = useCallback(async (newSettings: Settings) => {
    setLoading(true);
    setError(null);
    const result = await window.appAPI.user.saveSettings(newSettings);
    if (result.success) {
      setSettings(newSettings);
    } else {
      setError(result.error ?? null);
    }
    setLoading(false);
  }, []);

  return { settings, setSettings, loading, error, saveSettings };
};
