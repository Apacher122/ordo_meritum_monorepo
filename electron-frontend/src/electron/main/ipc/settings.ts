import { app, ipcMain, safeStorage, session } from "electron";

import fs from "fs";
import path from "path";

const settingsFilePath = path.join(app.getPath("userData"), "settings.json");

/**
 * Encrypts all string values in the provided apiKeys object
 * using the safeStorage module. If encryption fails for
 * any key, an error is logged to the console and the
 * original value is left unchanged.
 * @param {Record<string, any>} apiKeys - The object to
 *   encrypt. If the object is not provided or is not an
 *   object, this function does nothing.
 */
const encryptApiKeys = (apiKeys: Record<string, any> | undefined) => {
  if (!apiKeys || typeof apiKeys !== 'object') return;

  for (const key in apiKeys) {
    const val = apiKeys[key];
    if (typeof val === 'string' && val.trim().length > 0) {
      try {
        apiKeys[key] = safeStorage.encryptString(val).toString("base64");
      } catch (e) {
        console.error(`Encryption failed for ${key}:`, e);
      }
    }
  }
};

/**
 * Updates the User Agent string used by the default Electron session.
 * If the session does not exist or the userAgent string is not a string,
 * this function does nothing.
 * @param {string | undefined | null} userAgent - The User Agent
 *   string to set. If not a string, the default User Agent is used.
 */
const updateSessionUserAgent = (userAgent: unknown) => {
  if (session.defaultSession && typeof userAgent === 'string') {
    session.defaultSession.setUserAgent(userAgent);
  }
};

ipcMain.handle("save-settings", (_event, settings: any) => {
  try {
    const encrypted = structuredClone(settings);

    encryptApiKeys(encrypted.apiKeys);
    fs.writeFileSync(settingsFilePath, JSON.stringify(encrypted, null, 2));
    updateSessionUserAgent(settings.userAgent);

    return { success: true };
  } catch (error) {
    console.error("Failed to save settings:", error);
    return { 
      success: false, 
      error: error instanceof Error ? error.message : "Failed to save settings." 
    };
  }
});

ipcMain.handle("load-settings", () => {
  try {
    if (!fs.existsSync(settingsFilePath)) return { success: true, data: null };
    
    const rawData = fs.readFileSync(settingsFilePath, "utf-8");
    const data = JSON.parse(rawData);

    if (data.apiKeys && typeof data.apiKeys === 'object') {
      for (const key in data.apiKeys) {
        const enc = data.apiKeys[key];
        
        if (typeof enc === 'string' && enc.length > 0 && safeStorage.isEncryptionAvailable()) {
          try {
            const buf = Buffer.from(enc, "base64");
            data.apiKeys[key] = safeStorage.decryptString(buf);
          } catch (e) {
            console.error(`Decryption failed for ${key}:`, e);
            data.apiKeys[key] = "";
          }
        }
      }
    }
    
    return { success: true, data };
  } catch (error) {
    console.error("Failed to load settings:", error);
    return { success: false, error: "Failed to load settings." };
  }
});