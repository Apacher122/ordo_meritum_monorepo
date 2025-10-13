import { User, createUserWithEmailAndPassword, getAuth, signInWithEmailAndPassword } from 'firebase/auth';
import { doc, getDoc, getFirestore, setDoc } from 'firebase/firestore';
import { useEffect, useState } from 'react';

import { initializeApp } from 'firebase/app';

interface UserProfile {
  [key: string]: any;
}


/**
 * Hook to provide the Firebase Auth and Firestore instances,
 * as well as functions to register, log in, save and fetch user profiles.
 *
 * It uses the Firebase environment variables to initialize the app.
 *
 * @returns An object containing the `user`, `register`, `login`, `saveUserProfile`, and `fetchUserProfile` properties.
 */
export function useFirebaseAuth() {
  const [auth, setAuth] = useState<ReturnType<typeof getAuth> | null>(null);
  const [db, setDb] = useState<ReturnType<typeof getFirestore> | null>(null);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    if (!window.env) {
      console.error('Firebase env not found');
      return;
    }

    const firebaseConfig = {
      apiKey: window.env.FIREBASE_API_KEY,
      authDomain: window.env.FIREBASE_AUTH_DOMAIN,
      projectId: window.env.FIREBASE_PROJECT_ID,
      storageBucket: window.env.FIREBASE_STORAGE_BUCKET,
      messagingSenderId: window.env.FIREBASE_MESSAGING_SENDER_ID,
      appId: window.env.FIREBASE_APP_ID,
      measurementId: window.env.FIREBASE_MEASUREMENT_ID,
    };

    try {
      const app = initializeApp(firebaseConfig);
      setAuth(getAuth(app));
      setDb(getFirestore(app));
    } catch (err) {
      console.error('Firebase init failed:', err);
    }
  }, []);

  const register = async (email: string, password: string) => {
    if (!auth) throw new Error('Auth not initialized');
    const result = await createUserWithEmailAndPassword(auth, email, password);
    setUser(result.user);
    return result.user;
  };

  const login = async (email: string, password: string) => {
    if (!auth) throw new Error('Auth not initialized');
    const result = await signInWithEmailAndPassword(auth, email, password);
    setUser(result.user);
    return result.user;
  };

  const saveUserProfile = async (uid: string, profile: UserProfile) => {
    if (!db) throw new Error('Firestore not initialized');
    await setDoc(doc(db, 'users', uid), profile, { merge: true });
  };

  const fetchUserProfile = async (uid: string) => {
    if (!db) throw new Error('Firestore not initialized');
    const snap = await getDoc(doc(db, 'users', uid));
    return snap.exists() ? snap.data() : null;
  };

  return { user, register, login, saveUserProfile, fetchUserProfile };
}
