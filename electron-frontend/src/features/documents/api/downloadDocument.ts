import { ResumeChanges, CoverLetterChanges } from "../types";
import { apiRequest } from "@/shared/utils/requests";

const findInUint8Array = (haystack: Uint8Array, needle: Uint8Array): number => {
  for (let i = 0; i <= haystack.length - needle.length; i++) {
    let found = true;
    for (let j = 0; j < needle.length; j++) {
      if (haystack[i + j] !== needle[j]) {
        found = false;
        break;
      }
    }
    if (found) return i;
  }
  return -1;
};

const decoder = new TextDecoder("utf-8");
const uint8ArrayToString = (buffer: Uint8Array): string => {
  return decoder.decode(buffer);
};

export const downloadDocument = async (
  downloadUrl: string,
  changesUrl: string,
  token?: string
) => {
  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };

  const body = {
    download_url: downloadUrl,
    changes_url: changesUrl,
  };

  const multipartBlob = await apiRequest<Blob>('api/auth/downloads', { method: 'POST', headers, body, responseType: 'blob' });

  const arrayBuffer = await multipartBlob.arrayBuffer();
  const dataView = new Uint8Array(arrayBuffer);

  const firstChunkText = uint8ArrayToString(dataView.subarray(0, 256));
  const boundaryMatch = firstChunkText.match(/^--(.+)/);
  if (!boundaryMatch) {
    return { pdf: multipartBlob };
  }
  const boundary = `--${boundaryMatch[1]}`;
  const boundaryBytes = new TextEncoder().encode(boundary);

  let parts: Uint8Array[] = [];
  let start = findInUint8Array(dataView, boundaryBytes) + boundaryBytes.length;
  while (start !== -1 + boundaryBytes.length) {
    const end = findInUint8Array(dataView.subarray(start), boundaryBytes);
    if (end === -1) break;
    parts.push(dataView.subarray(start, start + end));
    start += end + boundaryBytes.length;
  }

  let pdf: Blob | null = null;
  let jsonData: ResumeChanges | CoverLetterChanges | undefined = undefined;

  const headerSeparator = new TextEncoder().encode("\r\n\r\n");

  for (const part of parts) {
    const headerEndIndex = findInUint8Array(part, headerSeparator);
    if (headerEndIndex === -1) continue;

    const headerText = uint8ArrayToString(part.subarray(2, headerEndIndex));

    const body = part.subarray(
      headerEndIndex + headerSeparator.length,
      part.length - 2
    );

    const filenameMatch = headerText.match(/filename="(.+?)"/);
    if (!filenameMatch) continue;
    const filename = filenameMatch[1];

    if (filename.endsWith(".pdf")) {
      pdf = new Blob([body.slice()], { type: "application/pdf" });
    } else if (filename.endsWith(".json")) {
      const jsonString = uint8ArrayToString(body);
      if (jsonString) {
        jsonData = JSON.parse(jsonString);
      }
    }
  }

  if (!pdf) {
    throw new Error("PDF file was not found in the server response.");
  }

  return { pdf, jsonData };
};
