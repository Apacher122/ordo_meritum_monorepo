import { ApiKeyConfig, RateLimitConfig, Settings } from "../types/types";
import { useCallback, useEffect, useState } from "react";

import { LlmProvider } from "@/shared/types";

const initialSettingsState: Settings = {
  apiKeys: {},
  featureAssignments: {
    matchSummary: { provider: "None", keyIndex: 0 },
    resumeGeneration: { provider: "None", keyIndex: 0 },
    coverLetterGeneration: { provider: "None", keyIndex: 0 },
  },
};

const isSamePTDay = (t1: number, t2: number) => {
  const fmt = new Intl.DateTimeFormat("en-US", {
    timeZone: "America/Los_Angeles",
    year: "numeric",
    month: "numeric",
    day: "numeric",
  });
  return fmt.format(new Date(t1)) === fmt.format(new Date(t2));
};

const updateRateLimit = (config: RateLimitConfig, now: number): RateLimitConfig => {
  let { currentDayCount, currentMinuteCount, lastDayReset, lastMinuteReset } = config;
  
  if (!isSamePTDay(lastDayReset, now)) {
    currentDayCount = 0;
    lastDayReset = now;
  }
  
  if (now - lastMinuteReset > 60000) {
    currentMinuteCount = 0;
    lastMinuteReset = now;
  }
  
  return { ...config, currentDayCount, currentMinuteCount, lastDayReset, lastMinuteReset };
};

const createDefaultRateLimit = (): RateLimitConfig => ({
  callsPerDay: 1000,
  callsPerMinute: 60,
  currentDayCount: 0,
  currentMinuteCount: 0,
  lastDayReset: Date.now(),
  lastMinuteReset: Date.now(),
});

/**
* Retrieves the application settings and provides functions to save and reset the settings.
*
* @returns {{
* settings: Settings | null,
* setSettings: (newSettings: Settings) => Promise<void>,
* loading: boolean,
* error: string | null,
* saveSettings: (newSettings: Settings) => Promise<void>,
* resetSettings: () => Promise<void>,
* getValidApiKey: (provider: LlmProvider) => Promise<string | null>,
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
        const loadedApiKeys = result.data.apiKeys || {};
        const migratedApiKeys: Partial<Record<LlmProvider, ApiKeyConfig[]>> =
        {};
        
        Object.keys(loadedApiKeys).forEach((key) => {
          const provider = key as LlmProvider;
          const value = loadedApiKeys[provider];
          
          if (Array.isArray(value)) {
            migratedApiKeys[provider] = value.map(
              (item: string | ApiKeyConfig) => {
                if (typeof item === "string") {
                  return {
                    key: item,
                    rateLimit: createDefaultRateLimit(),
                  };
                }
                return item;
              },
            );
          }
        });
        
        const mergedSettings: Settings = {
          apiKeys: {
            ...initialSettingsState.apiKeys,
            ...migratedApiKeys,
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
          cleanApiKeys[provider] = keys.map((k) => ({
            key: k.key,
            rateLimit: { ...k.rateLimit },
          }));
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
  
  const getValidApiKey = useCallback(
    async (provider: LlmProvider): Promise<string | null> => {
      const result = await window.appAPI.user.loadSettings();
      const currentSettings = result.success && result.data ? result.data : settings;
      
      if (!currentSettings?.apiKeys[provider]) return null;
      const keys = [...currentSettings.apiKeys[provider]];
      const now = Date.now();
      
      let selectedKeyIndex = -1;
      
      for (let i = 0; i < keys.length; i++) {
        const resetConfig = updateRateLimit(keys[i].rateLimit, now);
        
        if (
          resetConfig.currentDayCount < resetConfig.callsPerDay &&
          resetConfig.currentMinuteCount < resetConfig.callsPerMinute
        ) {
          selectedKeyIndex = i;
          keys[i].rateLimit = {
            ...resetConfig,
            currentDayCount: resetConfig.currentDayCount + 1,
            currentMinuteCount: resetConfig.currentMinuteCount + 1,
          };
          break;
        }
      }
      
      if (selectedKeyIndex !== -1) {
        const newSettings = {
          ...currentSettings,
          apiKeys: { ...currentSettings.apiKeys, [provider]: keys },
        };
        await window.appAPI.user.saveSettings(newSettings);
        setSettings(newSettings);
        return keys[selectedKeyIndex].key;
      }
      return null;
    },
    [settings]
  );
  
  return {
    settings,
    setSettings,
    loading,
    error,
    saveSettings,
    resetSettings,
    getValidApiKey,
  };
};
