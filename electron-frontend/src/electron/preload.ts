import { contextBridge, ipcRenderer } from "electron";

import { AppliedJob } from "@/features/applications/types";
import { Settings } from "@/features/settings/types/types";
import { UserProfile } from "@/features/user/types";

export interface WritingSample {
  fileName: string;
  content: string;
}

export interface IElectronAPI {
  user: {
    saveUserInfo: (
      userInfo: UserProfile
    ) => Promise<{ success: boolean; error?: string }>;
    loadUserInfo: () => Promise<{
      success: boolean;
      data?: UserProfile;
      error?: string;
    }>;
    saveSettings: (
      settings: Settings
    ) => Promise<{ success: boolean; error?: string }>;
    loadSettings: () => Promise<{
      success: boolean;
      data?: Settings;
      error?: string;
    }>;
  };
  files: {
    checkFileExists: (relativePath: string) => Promise<boolean>;
    saveFile: (
      relativePath: string,
      data: ArrayBuffer | string
    ) => Promise<{ success: boolean; path?: string; error?: string }>;
    saveJsonFile: (
      relativePath: string,
      data: any
    ) => Promise<{ success: boolean; error?: string }>;
    readJsonFile: (
      relativePath: string
    ) => Promise<{ success: boolean; data?: any; error?: string }>;
  };
  writingSamples: {
    upload: () => Promise<{
      success: boolean;
      samples?: WritingSample[];
      error?: string;
      reason?: string;
    }>;
    save: (
      samples: WritingSample[]
    ) => Promise<{ success: boolean; error?: string }>;
    load: () => Promise<{
      success: boolean;
      data?: WritingSample[];
      error?: string;
    }>;
  };
  polling: {
    start: () => void;
    stop: () => void;
    onUpdate: (callback: (jobs: AppliedJob[]) => void) => void;
  };
}

contextBridge.exposeInMainWorld("appAPI", {
  user: {
    saveUserInfo: (userInfo: UserProfile) =>
      ipcRenderer.invoke("save-user-info", userInfo),
    loadUserInfo: () => ipcRenderer.invoke("load-user-info"),
    saveSettings: (settings: Settings) =>
      ipcRenderer.invoke("save-settings", settings),
    loadSettings: () => ipcRenderer.invoke("load-settings"),
  },

  files: {
    checkFileExists: (relativePath: string) =>
      ipcRenderer.invoke("check-file-exists", relativePath),
    saveFile: (relativePath: string, data: ArrayBuffer | string) =>
      ipcRenderer.invoke("save-file", relativePath, data),
    saveJsonFile: (relativePath: string, data: any) =>
      ipcRenderer.invoke("save-json-file", relativePath, data),
    readJsonFile: (relativePath: string) =>
      ipcRenderer.invoke("read-json-file", relativePath),
  },

  writingSamples: {
    upload: () => ipcRenderer.invoke("upload-writing-samples"),
    save: (samples: WritingSample[]) =>
      ipcRenderer.invoke("save-writing-samples", samples),
    load: () => ipcRenderer.invoke("load-writing-samples"),
  },

  polling: {
    start: () => ipcRenderer.send("start-polling"),
    stop: () => ipcRenderer.send("stop-polling"),
    onUpdate: (callback: (jobs: AppliedJob[]) => void) =>
      ipcRenderer.on("jobs-updated", (_event, value) => callback(value)),
  },
});

contextBridge.exposeInMainWorld("env", {
  FIREBASE_API_KEY: process.env.FIREBASE_API_KEY,
  FIREBASE_AUTH_DOMAIN: process.env.FIREBASE_AUTH_DOMAIN,
  FIREBASE_PROJECT_ID: process.env.FIREBASE_PROJECT_ID,
  FIREBASE_STORAGE_BUCKET: process.env.FIREBASE_STORAGE_BUCKET,
  FIREBASE_MESSAGING_SENDER_ID: process.env.FIREBASE_MESSAGING_SENDER_ID,
  FIREBASE_APP_ID: process.env.FIREBASE_APP_ID,
  FIREBASE_MEASUREMENT_ID: process.env.FIREBASE_MEASUREMENT_ID,

  SERVER_URL: process.env.SERVER_URL,
});
