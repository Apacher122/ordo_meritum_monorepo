import { UserProfile } from "../types";

export const createSectionChangeHandler = <T extends keyof UserProfile>(
  section: T,
  setProfile: React.Dispatch<React.SetStateAction<UserProfile | null>>
) => {
  return (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setProfile((prev) => 
      prev && {
        ...prev,
        [section]: { ...(prev[section]), [name]: value },
      }
    );
  };
};