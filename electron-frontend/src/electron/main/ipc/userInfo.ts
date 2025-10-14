import { app, ipcMain } from "electron";

import fs from "fs";
import path from "path";

const userInfoFilePath = path.join(app.getPath("userData"), "userInfo.json");

ipcMain.handle("save-user-info", (event, userInfo) => {
  try {
    fs.writeFileSync(userInfoFilePath, JSON.stringify(userInfo, null, 2));
    return { success: true };
  } catch (error) {
    console.error("Failed to save user info:", error);
    return { success: false, error: "Failed to save data." };
  }
});

ipcMain.handle("load-user-info", () => {
  try {
    if (fs.existsSync(userInfoFilePath)) {
      const data = fs.readFileSync(userInfoFilePath, "utf-8");
      return { success: true, data: JSON.parse(data) };
    }
    return { success: true, data: null };
  } catch (error) {
    console.error("Failed to load user info:", error);
    return { success: false, error: "Failed to load data." };
  }
});
