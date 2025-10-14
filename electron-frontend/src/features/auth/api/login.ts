import { apiRequest } from "@/shared/utils/requests";

export const loginToServer = (): Promise<void> => {
  return apiRequest("/login-or-register", {
    method: "POST",
  });
};