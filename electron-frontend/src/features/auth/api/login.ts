import { apiRequest } from "@/shared/utils/requests";

/**
 * Sends a request to the backend to either log in an existing user or register
 * a new one based on the Firebase token provided in the authorization header.
 * @returns {Promise<void>} A promise that resolves upon a successful request.
 */
export const loginToServer = (): Promise<void> => {
  return apiRequest("/api/auth/login-or-register", {
    method: "POST",
  });
};