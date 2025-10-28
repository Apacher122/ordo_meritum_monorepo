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

/**
 * Provides a centralized state for managing the application's main header content,
 * including title, subtitle, and dynamic controls.
 * @param {object} props
 * @param {ReactNode} props.children - The child components that will have access to this context.
 * @returns {JSX.Element}
 */
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

/**
 * Custom hook to access the full HeaderContext.
 * @returns {HeaderContextType} The header context.
 * @throws {Error} If used outside of a HeaderProvider.
 */
export const useHeaderContext = () => {
  const context = useContext(HeaderContext);
  if (!context) {
    throw new Error("useHeaderContext must be used within a HeaderProvider");
  }
  return context;
};

/**
 * Custom hook for conveniently accessing the `setTitle` function from the HeaderContext.
 * @returns {(title: string) => void}
 */
export const useSetHeaderTitle = () => useHeaderContext().setTitle;

/**
 * Custom hook for conveniently accessing the `setSubtitle` function from the HeaderContext.
 * @returns {(subtitle: string) => void}
 */
export const useSetHeaderSubtitle = () => useHeaderContext().setSubtitle;

/**
 * Custom hook for conveniently accessing the `setControls` function from the HeaderContext.
 * @returns {(controls: ReactNode | null) => void}
 */
export const useSetHeaderControls = () => useHeaderContext().setControls;

