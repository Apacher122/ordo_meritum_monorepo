import { app, ipcMain } from "electron";

import fs from "fs";
import path from "path";
import sanitize from "sanitize-filename";

interface IpcResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  path?: string;
}

const basePublicDir = path.join(__dirname, "..", "public");
const paths = {
  pdfs: path.join(basePublicDir, "pdfs"),
  json: path.join(basePublicDir, "json"),
};

/**
 * Ensures that a directory exists, creating it recursively if it does not.
 * @param {string} dir - The directory path to ensure existence of.
 */
const ensureDir = (dir: string) => {
  if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
};

/**
 * Creates a safe, sanitized absolute file path.
 * @param {string} base - The base directory.
 * @param {string} relativePath - The user-provided relative path.
 * @returns {string} The sanitized, absolute file path.
 */
const getSafePath = (base: string, relativePath: string): string => {
  const safeName = sanitize(relativePath);
  return path.join(base, safeName);
};

/**
 * Writes data to a file at the specified absolute path.
 * @param {string} absolutePath - The absolute path to the file.
 * @param {string | Buffer} data - The data to write.
 */
const writeFile = (absolutePath: string, data: string | Buffer) => {
  ensureDir(path.dirname(absolutePath));
  fs.writeFileSync(absolutePath, data);
};

for (const dirPath of Object.values(paths)) {
  ensureDir(dirPath);
}

ipcMain.handle("read-pdf-bytes", async (_event, relativePath: string) => {
  try {
    const cleanPath = relativePath.replace("static://", "");
    const absolutePath = path.isAbsolute(cleanPath)
      ? cleanPath
      : path.join(app.getAppPath(), cleanPath);

    const buffer = await fs.readFileSync(absolutePath);
    return new Uint8Array(buffer);
  } catch (error) {
    console.error("IPC: Failed to read PDF file:", error);
    throw error;
  }
});

ipcMain.handle("check-file-exists", (event, relativePath: string): boolean => {
  const absPath = getSafePath(paths.pdfs, relativePath);
  return fs.existsSync(absPath);
});

ipcMain.handle(
  "save-file",
  (event, relativePath: string, data: ArrayBuffer | string): IpcResponse => {
    try {
      const absPath = getSafePath(paths.pdfs, relativePath);
      const bufferData = typeof data === "string" ? data : Buffer.from(data);
      writeFile(absPath, bufferData);
      return { success: true, path: `pdfs/${relativePath}` };
    } catch (error) {
      console.error(`Failed to save file (${relativePath}):`, error);
      return { success: false, error: "Failed to save file." };
    }
  },
);

ipcMain.handle(
  "save-json-file",
  (event, relativePath: string, data: any): IpcResponse => {
    try {
      const absPath = getSafePath(paths.json, relativePath);
      writeFile(absPath, JSON.stringify(data, null, 2));
      return { success: true };
    } catch (error) {
      console.error(`Failed to save JSON file (${relativePath}):`, error);
      return { success: false, error: "Failed to save JSON file." };
    }
  },
);

ipcMain.handle(
  "read-json-file",
  (event, relativePath: string): IpcResponse<any> => {
    try {
      const absPath = getSafePath(paths.json, relativePath);
      if (!fs.existsSync(absPath)) {
        return { success: false, error: "File not found." };
      }
      const fileContent = fs.readFileSync(absPath, "utf-8");
      return { success: true, data: JSON.parse(fileContent) };
    } catch (error) {
      console.error(`Failed to read JSON file (${relativePath}):`, error);
      return { success: false, error: "Failed to read or parse JSON file." };
    }
  },
);
