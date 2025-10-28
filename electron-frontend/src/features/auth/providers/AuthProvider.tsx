import React, {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useMemo,
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
 * Provides authentication state to its children components. It listens for
 * Firebase auth state changes and syncs the user with the backend upon login.
 * @param {object} props - The component props.
 * @param {ReactNode} props.children - The child components to render.
 * @returns {JSX.Element}
 */
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const value = useMemo(() => ({ user, loading }), [user, loading]);

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
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

/**
 * Custom hook to access the authentication context.
 * @returns {AuthContextType} The authentication context, including the user and loading state.
 * @throws {Error} If used outside of an `AuthProvider`.
 */
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};