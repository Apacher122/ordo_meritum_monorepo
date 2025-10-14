import { app, dialog, ipcMain } from "electron";

import fs from "fs";
import { getMainWindow } from "../window";
import path from "path";

const writingSamplesPath = path.join(
  app.getPath("userData"),
  "writingSamples.json"
);

ipcMain.handle("upload-writing-samples", async () => {
  const mainWindow = getMainWindow();
  if (!mainWindow) {
    return { success: false, error: "Main window not available." };
  }

  const result = (await dialog.showOpenDialog(mainWindow, {
    title: "Select Writing Samples",
    properties: ["openFile", "multiSelections"],
    filters: [{ name: "Text Files", extensions: ["txt"] }],
  })) as unknown as { canceled: boolean; filePaths: string[] };

  if (result.canceled || result.filePaths.length === 0) {
    return { success: false, reason: "Dialog canceled." };
  }

  try {
    const samples = result.filePaths.map((filePath: string) => ({
      fileName: path.basename(filePath),
      content: fs.readFileSync(filePath, "utf-8"),
    }));
    return { success: true, samples };
  } catch (error) {
    console.error("Failed to read writing samples:", error);
    return { success: false, error: "Failed to read files." };
  }
});

ipcMain.handle("save-writing-samples", (_event, samples) => {
  try {
    fs.writeFileSync(writingSamplesPath, JSON.stringify(samples, null, 2));
    return { success: true };
  } catch (error) {
    console.error("Failed to save writing samples:", error);
    return { success: false, error: "Failed to save samples." };
  }
});

ipcMain.handle("load-writing-samples", () => {
  try {
    if (fs.existsSync(writingSamplesPath)) {
      const data = fs.readFileSync(writingSamplesPath, "utf-8");
      return { success: true, data: JSON.parse(data) };
    }
    return { success: true, data: [] };
  } catch (error) {
    console.error("Failed to load writing samples:", error);
    return { success: false, error: "Failed to load samples." };
  }
});
