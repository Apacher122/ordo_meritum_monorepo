import * as path from "path";

export const validatePath = (filePath: string): string => {
  if (path.isAbsolute(filePath) && !filePath.includes("..")) {
    return filePath;
  }
  throw new Error("Invalid file path");
};
