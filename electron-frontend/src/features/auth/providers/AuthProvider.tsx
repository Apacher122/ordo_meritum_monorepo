import React, {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { User, onAuthStateChanged } from "firebase/auth";

import { auth } from "@/config/firebase";
import { loginToServer } from "../api/login";

interface AuthContextType {
  user: User | null;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * A React context provider that provides the state and functions
 * for managing the user state and sync with the backend.
 *
 * The provider wraps the `useAuth` hook and provides
 * the following values to its children:
 *
 *   - `user`: The currently signed-in user.
 *   - `loading`: A boolean indicating whether the user is currently being fetched.
 */
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
      setUser(currentUser);
      setLoading(false);
    });
    return () => unsubscribe();
  }, []);

  useEffect(() => {
    const syncUserWithBackend = async () => {
      if (user) {
        try {
          console.log("User logged in, syncing with backend...");
          await loginToServer();
          console.log("Backend sync successful.");
        } catch (error) {
          console.error("Failed to sync user with backend:", error);
        }
      }
    };

    if (!loading) {
      syncUserWithBackend();
    }
  }, [user, loading]); 

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100">
        <div>Loading Application...</div>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={{ user, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};