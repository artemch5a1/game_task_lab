import type { JSX } from "solid-js";
import "./Layout.css";
import { Sidebar } from "../shared/components/sidebar/Sidebar.tsx";

interface AppLayoutProps {
  children: JSX.Element;
  showSidebar?: boolean;
}

export const AppLayout = (props: AppLayoutProps) => {
  const showSidebar = props.showSidebar ?? true;

  return (
    <div class="app-shell">
      {showSidebar && <Sidebar />}
      <main class="app-main app-content">{props.children}</main>
    </div>
  );
};

