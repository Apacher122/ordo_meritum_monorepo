import { useEffect, useState } from "react";

import { User } from "@/shared/types";
import { getPublicKey } from "../api/getPublicKey";

/**
 * A custom hook to get the public key and listen for real-time updates via SSE.
 * It will only run if a user is provided.
 * @param {User | null} user - The authenticated user object.
 * @returns An object containing the public key, loading state, and any error.
 */
export const usePublicKeyStream = (user: User | null) => {
  const [publicKey, setPublicKey] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) {
      setLoading(false);
      setPublicKey(null);
      setError(null);
      return;
    }

    setLoading(true);

    const fetchInitialKey = async () => {
      try {
        const key = await getPublicKey();
        setPublicKey(key);
      } catch (err) {
        console.error("Failed to fetch public key:", err);
        setError("Failed to load public key.");
      } finally {
        setLoading(false);
      }
    };

    fetchInitialKey();

    const eventSource = new EventSource(`${window.env.SERVER_URL}/public-key-stream`);

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.key) {
        setPublicKey(data.key);
      }
    };

    eventSource.onerror = (err) => {
      console.error("EventSource failed:", err);
      eventSource.close();
      setError("Connection to key stream lost.");
    };

    return () => {
      eventSource.close();
    };
  }, [user]);

  return { publicKey, loading, error };
};
