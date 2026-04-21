import { app } from "electron";
import fs from "fs";
import path from "path";

export const getFilePath = (fileName: string) => path.join(app.getPath("userData"), fileName);

export const loadJsonFile = (fileName: string) => {
  const filePath = getFilePath(fileName);
  try {
    if (!fs.existsSync(filePath)) return { success: true, data: null };
    const data = fs.readFileSync(filePath, "utf-8");
    return { success: true, data: JSON.parse(data) };
  } catch (error) {
    console.error(`Failed to load ${fileName}:`, error);
    return { success: false, error: `Failed to load ${fileName}.` };
  }
};

export const saveJsonFile = (fileName: string, data: any) => {
  try {
    const filePath = getFilePath(fileName);
    fs.writeFileSync(filePath, JSON.stringify(data, null, 2));
    return { success: true };
  } catch (error) {
    console.error(`Failed to save ${fileName}:`, error);
    return { success: false, error: `Failed to save ${fileName}.` };
  }
};