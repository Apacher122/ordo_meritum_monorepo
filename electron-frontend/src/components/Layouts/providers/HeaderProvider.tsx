import React, {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState
} from "react";

interface HeaderState {
  title: string;
  subtitle: string;
  controls: ReactNode | null;
}

interface HeaderContextType extends HeaderState {
  setTitle: (title: string) => void;
  setSubtitle: (subtitle: string) => void;
  setControls: (controls: ReactNode | null) => void;
}

const HeaderContext = createContext<HeaderContextType | undefined>(undefined);

export const HeaderProvider = ({ children }: { children: ReactNode }) => {
  const [headerState, setHeaderState] = useState<HeaderState>({
    title: "No Job Selected",
    subtitle: "Select or analyze a job to begin",
    controls: null,
  });

  const setTitle = useCallback((title: string) => setHeaderState(s => ({ ...s, title })), []);
  const setSubtitle = useCallback((subtitle: string) => setHeaderState(s => ({ ...s, subtitle })), []);
  const setControls = useCallback((controls: ReactNode | null) => setHeaderState(s => ({ ...s, controls })), []);

  const value = useMemo(
    () => ({
      ...headerState,
      setTitle,
      setSubtitle,
      setControls,
    }),
    [headerState, setTitle, setSubtitle, setControls]
  );

  return (
    <HeaderContext.Provider value={value}>{children}</HeaderContext.Provider>
  );
};

export const useHeaderContext = () => {
  const context = useContext(HeaderContext);
  if (!context) {
    throw new Error("useHeaderContext must be used within a HeaderProvider");
  }
  return context;
};

export const useSetHeaderTitle = () => useHeaderContext().setTitle;
export const useSetHeaderSubtitle = () => useHeaderContext().setSubtitle;
export const useSetHeaderControls = () => useHeaderContext().setControls;

