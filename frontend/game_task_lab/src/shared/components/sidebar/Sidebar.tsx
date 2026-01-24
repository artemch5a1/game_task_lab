import "../../../assets/Layout.css";
import { useNavigate } from "@solidjs/router";
import { authStore } from "../../../features/auth/store/auth.store";

export const Sidebar = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    authStore.actions.logout();
    navigate("/login", { replace: true });
  };

  return (
    <aside class="app-sidebar">
      <div class="app-sidebar-header">
        <span class="app-logo-mark">G</span>
        <span class="app-logo-text">GameLab</span>
      </div>

      <nav class="app-sidebar-section">
        <ul class="app-sidebar-nav">
          <li class="app-sidebar-item app-sidebar-item--active">
            <span class="app-sidebar-label">GameList</span>
          </li>
        </ul>
      </nav>

      <div class="app-sidebar-footer">
        <button type="button" class="app-sidebar-footer-item" onClick={handleLogout}>
          <span class="app-sidebar-icon">⎋</span>
          <span class="app-sidebar-label">Выйти</span>
        </button>
      </div>
    </aside>
  );
};

