import { BrowserWindow } from "electron";
import path from "path";

export let mainWindow: BrowserWindow | null = null;

export const createWindow = () => {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      preload: path.join(__dirname, "preload.js"),
    },
  });

  if (process.env.NODE_ENV === "development") {
    const hostUrl = process.env.SERVER_URL ? process.env.SERVER_URL : "http://localhost:8080";
    mainWindow.loadURL(hostUrl);
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadFile(path.join(__dirname, "index.html"));
  }

  mainWindow.on("closed", () => (mainWindow = null));
};
