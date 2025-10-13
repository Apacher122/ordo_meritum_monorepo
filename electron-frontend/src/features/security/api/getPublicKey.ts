import { apiRequest } from "@/shared/utils/requests";

/**
 * Fetches the public key from the standard JSON endpoint.
 * @returns {Promise<string>} A promise that resolves to the public key string.
 */
export const getPublicKey = async (): Promise<string> => {
  console.log("Fetching public key...");
  const response = await apiRequest<{ publicKey: string }>("/public-key");
  console.log(`Received public key: ${response.publicKey}`);
  return response.publicKey;
};
