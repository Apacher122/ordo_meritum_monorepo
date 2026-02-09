import "./ipc"; // auto-loads all IPC handlers

import { BrowserWindow, app } from "electron";

import { createWindow } from "./window";
import { initWebSocketIpc } from "./ipc/websocket";

app.whenReady().then(() => {
  const mainWindow = createWindow();
  if (mainWindow) {
    initWebSocketIpc(mainWindow);
  }
  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") app.quit();
});
