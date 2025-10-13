import "../../assets/styles/components/MainShell.css";

import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { getAuth, signOut } from "firebase/auth";

import { Footer } from "./Footer";
import { Header } from "./Header";
import React from "react";
import Sidebar from "./Sidebar";

const SIDEBAR_ITEMS: { id: string; label: string }[] = [
  { id: "/info", label: "Job Info" },
  { id: "/match-summary", label: "Match Summary" },
  { id: "/applications", label: "Applications" },
  { id: "/documents", label: "Documents" },
  { id: "/user-info", label: "User Profile" },
  { id: "/settings", label: "Settings" },
  { id: "signout", label: "Sign Out" },
];

export const MainShell: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const handleSidebarSelect = async (id: string) => {
    if (id === "signout") {
      await signOut(getAuth());
      navigate("/login");
      return;
    }
    navigate(id);
  };

  return (
    <div className="main-shell">
      <Sidebar
        items={SIDEBAR_ITEMS}
        activeId={location.pathname}
        onSelect={handleSidebarSelect}
      />
      <div className="content-wrapper">
        <Header /> 
        <main className="page-content">
          <Outlet />
        </main>
        <Footer />
      </div>
    </div>
  );
};