import { apiRequest } from "@/shared/utils/requests";

export const loginToServer = (): Promise<void> => {
  return apiRequest("api/secure/auth/login", {
    method: "POST",
  });
};