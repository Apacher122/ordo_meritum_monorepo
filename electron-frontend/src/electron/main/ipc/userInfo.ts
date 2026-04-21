import { loadJsonFile, saveJsonFile } from "@/shared/utils/fileManager";

import { ipcMain } from "electron";

const FILE_NAME = "userInfo.json";

/**
 * Handles the "save-user-info" IPC event. Saves the user's profile data
 * to a local JSON file.
 */
ipcMain.handle("save-user-info", (event, userInfo) => saveJsonFile(FILE_NAME, userInfo));

/**
 * Handles the "load-user-info" IPC event. Loads the user's profile data
 * from a local JSON file.
 */
ipcMain.handle("load-user-info", () => loadJsonFile(FILE_NAME));