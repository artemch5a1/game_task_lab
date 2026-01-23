import type { JSX } from "solid-js";
import "./Layout.css";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";

interface AppLayoutProps {
  children: JSX.Element;
}

export const AppLayout = (props: AppLayoutProps) => {
  return (
    <div class="app-shell">
      <Sidebar />

      <div class="app-main">
        <Topbar />
        <main class="app-content">{props.children}</main>
      </div>
    </div>
  );
};

