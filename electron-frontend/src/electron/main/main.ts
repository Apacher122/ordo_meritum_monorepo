import "./ipc"; // auto-loads all IPC handlers

import { BrowserWindow, app, net, protocol } from "electron";

import { createWindow } from "./window";
import { initWebSocketIpc } from "./ipc/websocket";
import path from "path";
import { pathToFileURL } from "url";

protocol.registerSchemesAsPrivileged([
  {
    scheme: "static",
    privileges: {
      standard: true,
      secure: true,
      supportFetchAPI: true,
      bypassCSP: false,
    },
  },
]);

app.whenReady().then(() => {
  protocol.handle("static", (request) => {
    const url = request.url.replace("static://", "");
    const decodedPath = decodeURIComponent(url);
    const absolutePath = path.isAbsolute(decodedPath)
      ? decodedPath
      : path.join(app.getAppPath(), decodedPath);

    return net.fetch(pathToFileURL(absolutePath).toString());
  });
  
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
