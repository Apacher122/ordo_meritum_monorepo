import { auth } from "@/config/firebase";

const API_BASE_URL = "http://localhost:8080";

interface RequestOptions extends RequestInit {
  body?: any;
  responseType?: "json" | "text" | "blob";
}

/**
 * A centralized and authenticated API request helper function.
 * It automatically adds the Firebase Auth ID token to the Authorization header
 * for every request, ensuring all outgoing calls are secure.
 *
 * @param {string} endpoint - The API endpoint to call (e.g., 'user/application-list').
 * @param {RequestOptions} [options={}] - Optional fetch options (method, body, headers, etc.).
 * @returns {Promise<T>} A promise that resolves to the parsed response data.
 * @template T - The expected type of the response data.
 */
export const apiRequest = async <T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> => {
  const { body, responseType = "json", ...customConfig } = options;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (customConfig.headers) {
    Object.assign(headers, customConfig.headers);
  }

  if (auth.currentUser) {
    try {
      const token = await auth.currentUser.getIdToken();
      headers["Authorization"] = `Bearer ${token}`;
    } catch (error) {
      console.error("Failed to get auth token:", error);
    }
  }

  const config: RequestInit = {
    method: body ? "POST" : "GET",
    ...customConfig,
    headers,
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  const response = await fetch(`${window.env.SERVER_URL}/${endpoint}`, config);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({
      message: "An unknown API error occurred.",
    }));
    throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
  }

  if (responseType === "blob") {
    return response.blob() as Promise<T>;
  }
  if (responseType === "text") {
    return response.text() as Promise<T>;
  }
  return response.json() as Promise<T>;
};