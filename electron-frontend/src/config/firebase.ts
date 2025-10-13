import {
  Auth,
  GoogleAuthProvider,
  browserLocalPersistence,
  getAuth,
  setPersistence,
} from "firebase/auth";
import { FirebaseApp, initializeApp } from "firebase/app";
import { Firestore, getFirestore } from "firebase/firestore";

const firebaseConfig = {
  apiKey: window.env.FIREBASE_API_KEY ?? "",
  authDomain: window.env.FIREBASE_AUTH_DOMAIN ?? "",
  projectId: window.env.FIREBASE_PROJECT_ID ?? "",
  storageBucket: window.env.FIREBASE_STORAGE_BUCKET ?? "",
  messagingSenderId: window.env.FIREBASE_MESSAGING_SENDER_ID ?? "",
  appId: window.env.FIREBASE_APP_ID ?? "",
  measurementId: window.env.FIREBASE_MEASUREMENT_ID ?? "",
};

export const app: FirebaseApp = initializeApp(firebaseConfig);

export const auth: Auth = getAuth(app);
export const db: Firestore = getFirestore(app);
export const googleProvider = new GoogleAuthProvider();

setPersistence(auth, browserLocalPersistence).catch((error) => {
  console.error("Firebase persistence error:", error);
});
