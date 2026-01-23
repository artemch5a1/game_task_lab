import type { JSX } from "solid-js";
import "./Layout.css";
import { Sidebar } from "../shared/components/sidebar/Sidebar.tsx";

interface AppLayoutProps {
  children: JSX.Element;
}

export const AppLayout = (props: AppLayoutProps) => {
  return (
    <div class="app-shell">
      <Sidebar />
      <main class="app-main app-content">{props.children}</main>
    </div>
  );
};

