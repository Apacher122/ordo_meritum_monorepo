import { app, ipcMain, safeStorage } from "electron";

import fs from "fs";
import path from "path";

const settingsFilePath = path.join(app.getPath("userData"), "settings.json");

ipcMain.handle("save-settings", (event, settings) => {
  try {
    const encrypted = { ...settings };
    if (encrypted.apiKeys) {
      for (const key in encrypted.apiKeys) {
        const val = encrypted.apiKeys[key];
        if (val) {
          encrypted.apiKeys[key] = safeStorage.encryptString(val).toString("base64");
        }
      }
    }
    fs.writeFileSync(settingsFilePath, JSON.stringify(encrypted, null, 2));
    return { success: true };
  } catch (error) {
    console.error("Failed to save settings:", error);
    return { success: false, error: "Failed to save settings." };
  }
});

ipcMain.handle("load-settings", () => {
  try {
    if (!fs.existsSync(settingsFilePath)) return { success: true, data: null };
    const data = JSON.parse(fs.readFileSync(settingsFilePath, "utf-8"));
    if (data.apiKeys) {
      for (const key in data.apiKeys) {
        const enc = data.apiKeys[key];
        if (enc && safeStorage.isEncryptionAvailable()) {
          const buf = Buffer.from(enc, "base64");
          data.apiKeys[key] = safeStorage.decryptString(buf);
        }
      }
    }
    return { success: true, data };
  } catch (error) {
    console.error("Failed to load settings:", error);
    return { success: false, error: "Failed to load settings." };
  }
});
