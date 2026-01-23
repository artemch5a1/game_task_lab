import "./Layout.css";

const navItemsTop = [
  { label: "Dashboard", icon: "🏠" },
  { label: "Products", icon: "📦" },
  { label: "Favorites", icon: "❤️" },
  { label: "Inbox", icon: "📥", active: true },
  { label: "Order Lists", icon: "🧾" },
  { label: "Product Stock", icon: "📊" },
];

const navItemsPages = [
  "Pricing",
  "Calendar",
  "To-Do",
  "Contact",
  "Invoice",
  "UI Elements",
  "Team",
  "Table",
];

const navItemsBottom = [
  { label: "Settings", icon: "⚙️" },
  { label: "Logout", icon: "⏻" },
];

export const Sidebar = () => {
  return (
    <aside class="app-sidebar">
      <div class="app-sidebar-header">
        <span class="app-logo-mark">D</span>
        <span class="app-logo-text">DashStack</span>
      </div>

      <nav class="app-sidebar-section">
        <div class="app-sidebar-section-title">MENU</div>
        <ul class="app-sidebar-nav">
          {navItemsTop.map((item) => (
            <li
              class={`app-sidebar-item ${
                item.active ? "app-sidebar-item--active" : ""
              }`}
            >
              <span class="app-sidebar-icon">{item.icon}</span>
              <span class="app-sidebar-label">{item.label}</span>
            </li>
          ))}
        </ul>
      </nav>

      <nav class="app-sidebar-section">
        <div class="app-sidebar-section-title">PAGES</div>
        <ul class="app-sidebar-nav">
          {navItemsPages.map((label) => (
            <li class="app-sidebar-item app-sidebar-item--sub">
              <span class="app-sidebar-label">{label}</span>
            </li>
          ))}
        </ul>
      </nav>

      <div class="app-sidebar-footer">
        {navItemsBottom.map((item) => (
          <button class="app-sidebar-footer-item" type="button">
            <span class="app-sidebar-icon">{item.icon}</span>
            <span class="app-sidebar-label">{item.label}</span>
          </button>
        ))}
      </div>
    </aside>
  );
};

