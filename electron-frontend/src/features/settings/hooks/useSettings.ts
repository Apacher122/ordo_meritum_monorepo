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


/**
 * Hook to manage the application settings.
 * It loads the settings from the electron storage on mount and provides functions to update the settings.
 * It also provides the current state of the settings, whether the settings are being loaded, and any errors that may have occurred.
 * @returns {{ settings: Settings, setSettings: (newSettings: Settings) => Promise<void>, loading: boolean, error: string | null, saveSettings: (newSettings: Settings) => Promise<void> }}
 */

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

        if (!loadedData) {
          setSettings(initialSettingsState);
        } else {
          const mergedSettings: Settings = {
            apiKeys: {
              ...initialSettingsState.apiKeys,
              ...(loadedData.apiKeys || {}),
            },
            featureAssignments: {
              ...initialSettingsState.featureAssignments,
              ...(loadedData.featureAssignments || {}),
            },
          };
          setSettings(mergedSettings);
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
