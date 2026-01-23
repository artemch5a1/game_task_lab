import "../../../assets/Layout.css";

export const Sidebar = () => {
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
    </aside>
  );
};

