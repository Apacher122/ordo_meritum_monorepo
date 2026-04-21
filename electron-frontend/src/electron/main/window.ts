import { BrowserWindow } from "electron";
import path from "path";

let mainWindow: BrowserWindow | null = null;

/**
 * Creates and configures the main application window.
 */
export const createWindow = () => {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      preload: path.join(__dirname, "preload.js"),
      sandbox: false,
      webSecurity: false,
    },
  });

  if (process.env.NODE_ENV === "development") {
    const hostUrl = process.env.SERVER_URL
      ? process.env.SERVER_URL
      : "http://localhost:8080";
    mainWindow.loadURL(hostUrl);
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadFile(path.join(__dirname, "index.html"));
  }

  mainWindow.on("closed", () => (mainWindow = null));

  return mainWindow;
};

/**
 * Returns the main application window instance.
 * @returns {BrowserWindow | null} The main window instance, or null if it has been closed.
 */
export const getMainWindow = (): BrowserWindow | null => {
  return mainWindow;
};
