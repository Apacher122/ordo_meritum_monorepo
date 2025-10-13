import * as providers from "./appProviders";

import React from "react";
import { HashRouter as Router } from "react-router-dom";

export const AppProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <Router>
      <providers.AuthProvider>
        <providers.ApplicationProvider>
          <providers.DocumentStatusProvider>
            <providers.HeaderProvider>
              <providers.PublicKeyProvider>
              {children}
              </providers.PublicKeyProvider>
            </providers.HeaderProvider>
          </providers.DocumentStatusProvider>
        </providers.ApplicationProvider>
      </providers.AuthProvider>
    </Router>
  );
};