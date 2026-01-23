import "./Layout.css";

export const Topbar = () => {
  return (
    <header class="app-topbar">
      <button class="app-topbar-menu-button" type="button">
        ☰
      </button>

      <div class="app-topbar-search">
        <span class="app-topbar-search-icon">🔍</span>
        <input
          class="app-topbar-search-input"
          type="text"
          placeholder="Search"
        />
      </div>

      <div class="app-topbar-right">
        <button class="app-topbar-icon-button" type="button">
          🔔
          <span class="app-topbar-badge">6</span>
        </button>

        <div class="app-topbar-language">
          <span role="img" aria-label="flag">
            🇬🇧
          </span>
          <span class="app-topbar-language-label">English</span>
          <span class="app-topbar-language-caret">▾</span>
        </div>

        <div class="app-topbar-user">
          <div class="app-topbar-avatar">MR</div>
          <div class="app-topbar-user-info">
            <div class="app-topbar-user-name">Moni Roy</div>
            <div class="app-topbar-user-role">Admin</div>
          </div>
        </div>
      </div>
    </header>
  );
};

