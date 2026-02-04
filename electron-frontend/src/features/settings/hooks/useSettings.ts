import { useCallback, useEffect, useState } from "react";

import { Settings } from "../types/types";

/**
 * Default state for the application settings.
 */
const initialSettingsState: Settings = {
  apiKeys: {},
  featureAssignments: {
    matchSummary: { provider: "None", keyIndex: 0 },
    resumeGeneration: { provider: "None", keyIndex: 0 },
    coverLetterGeneration: { provider: "None", keyIndex: 0 },
  },
};

/**
 * Retrieves the application settings and provides functions to save and reset the settings.
 *
 * @returns {{
 *   settings: Settings | null,
 *   setSettings: (newSettings: Settings) => Promise<void>,
 *   loading: boolean,
 *   error: string | null,
 *   saveSettings: (newSettings: Settings) => Promise<void>,
 *   resetSettings: () => Promise<void>,
 * }}
 */
export const useSettings = () => {
  const [settings, setSettings] = useState<Settings | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await window.appAPI.user.loadSettings();
      if (result.success && result.data) {
        const mergedSettings: Settings = {
          apiKeys: {
            ...initialSettingsState.apiKeys,
            ...result.data.apiKeys,
          },
          featureAssignments: {
            ...initialSettingsState.featureAssignments,
            ...result.data.featureAssignments,
          },
        };
        setSettings(mergedSettings);
      } else {
        setSettings(initialSettingsState);
      }
    } catch (err) {
      console.error("Failed to load settings:", err);
      setError("Failed to initialize settings.");
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);
  const saveSettings = useCallback(async (newSettings: Settings) => {
    if (!newSettings) return;

    setLoading(true);
    setError(null);
    try {
      const cleanApiKeys: any = {};
      Object.entries(newSettings.apiKeys || {}).forEach(([provider, keys]) => {
        if (Array.isArray(keys)) {
          cleanApiKeys[provider] = [...keys];
        }
      });

      const cleanAssignments: any = {};
      Object.entries(newSettings.featureAssignments || {}).forEach(
        ([feature, assignment]) => {
          cleanAssignments[feature] = {
            provider: assignment.provider,
            keyIndex: Number(assignment.keyIndex) || 0,
          };
        },
      );

      const finalPayload = {
        apiKeys: cleanApiKeys,
        featureAssignments: cleanAssignments,
      };

      const result = await window.appAPI.user.saveSettings(finalPayload);

      if (result.success) {
        setSettings(newSettings);
      } else {
        setError(result.error ?? "Failed to save settings.");
      }
    } catch (err: any) {
      console.error("IPC Save Error:", err);
      setError(`Save Error: ${err.message || "Object conversion failed"}`);
    }
    setLoading(false);
  }, []);

  const resetSettings = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const cleanDefaults: Settings = {
        apiKeys: {},
        featureAssignments: {
          matchSummary: { provider: "None", keyIndex: 0 },
          resumeGeneration: { provider: "None", keyIndex: 0 },
          coverLetterGeneration: { provider: "None", keyIndex: 0 },
        },
      };

      const result = await window.appAPI.user.saveSettings(cleanDefaults);
      if (result.success) {
        setSettings(initialSettingsState);
      } else {
        setError(result.error ?? "Failed to reset settings.");
      }
    } catch (err: any) {
      setError("Failed to clear existing data.");
      console.error("IPC Reset Error:", err);
    }
    setLoading(false);
  }, []);

  return { settings, setSettings, loading, error, saveSettings, resetSettings };
};
