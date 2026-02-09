import WebSocket from "ws";
import { ipcMain } from "electron";

let socket: WebSocket | null = null;

/**
 * Initializes IPC handlers for WebSocket connections.
 * Listens for "ws-connect" event and sets up a WebSocket connection
 * with the provided URL.
 * Listens for "ws-send" event and sends a message over the WebSocket connection.
 * Listens for "ws-disconnect" event and closes the WebSocket connection.
 * @param {Electron.BrowserWindow} mainWindow - The main window of the Electron application.
 */
export const initWebSocketIpc = (mainWindow: Electron.BrowserWindow) => {
  ipcMain.on("ws-connect", (event, { url }) => {
    if (socket) {
      socket.terminate();
    }

    console.log("Main Process: Connecting to", url);
    socket = new WebSocket(url);

    socket.on("open", () => {
      mainWindow.webContents.send("ws-status", "connected");
    });

    socket.on("message", (data) => {
      const messageString =
        typeof data === "string" ? data : JSON.stringify(data);

      console.log("Main Process: Received message", messageString);
      mainWindow.webContents.send("ws-message", messageString);
    });

    socket.on("error", (err) => {
      console.error("Main Process WS Error:", err.message);
      mainWindow.webContents.send("ws-error", err.message);
    });

    socket.on("close", () => {
      mainWindow.webContents.send("ws-status", "disconnected");
    });
  });

  ipcMain.on("ws-send", (event, message) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(message);
    }
  });

  ipcMain.on("ws-disconnect", () => {
    socket?.close();
    socket = null;
  });
};
