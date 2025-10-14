import { getMainWindow } from "@/electron/main/";
import { ipcMain } from "electron";

let pollingInterval: NodeJS.Timeout | null = null;

async function fetchJobsFromBackend() {
  console.log("Polling for new jobs...");
  return [];
}

ipcMain.on("start-polling", () => {
  const mainWindow = getMainWindow();
  if (pollingInterval) clearInterval(pollingInterval);
  pollingInterval = setInterval(async () => {
    if (mainWindow) {
      const jobs = await fetchJobsFromBackend();
      mainWindow.webContents.send("jobs-updated", jobs);
    }
  }, 5000);
});

ipcMain.on("stop-polling", () => {
  if (pollingInterval) {
    clearInterval(pollingInterval);
    pollingInterval = null;
  }
});