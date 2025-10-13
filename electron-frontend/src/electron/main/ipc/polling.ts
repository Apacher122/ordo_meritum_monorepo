import { ipcMain } from "electron";
import { mainWindow } from "@/electron/main/";


let pollingInterval: NodeJS.Timeout | null = null;

async function fetchJobsFromBackend() {
  console.log("Polling for new jobs...");
  return [];
}

ipcMain.on("start-polling", () => {
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