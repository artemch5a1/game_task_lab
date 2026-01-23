import "./GamesPage.css";
import { GameList } from "../components/GameList";

export const GamesPage = () => {
  return (
    <div class="games-page">
      <section class="games-page-header">
        <div>
          <h1 class="games-page-title">Games</h1>
          <p class="games-page-subtitle">Manage and explore your game library</p>
        </div>
        <button class="games-page-primary-button" type="button">
          + New Game
        </button>
      </section>

      <section class="games-page-body">
        <div class="games-page-sidebar-card">
          <button class="games-page-compose-button" type="button">
            + Compose
          </button>

          <div class="games-page-sidebar-block">
            <div class="games-page-sidebar-block-title">My Lists</div>
            <ul class="games-page-sidebar-list">
              <li class="games-page-sidebar-item games-page-sidebar-item--active">
                <span>All games</span>
              </li>
              <li class="games-page-sidebar-item">
                <span>Favorites</span>
              </li>
              <li class="games-page-sidebar-item">
                <span>Recently added</span>
              </li>
            </ul>
          </div>

          <div class="games-page-sidebar-block">
            <div class="games-page-sidebar-block-title">Labels</div>
            <ul class="games-page-sidebar-labels">
              <li>
                <span class="games-page-label-dot games-page-label-dot--primary" />
                Primary
              </li>
              <li>
                <span class="games-page-label-dot games-page-label-dot--social" />
                Social
              </li>
              <li>
                <span class="games-page-label-dot games-page-label-dot--work" />
                Work
              </li>
              <li>
                <span class="games-page-label-dot games-page-label-dot--friends" />
                Friends
              </li>
            </ul>
          </div>
        </div>

        <div class="games-page-main-card">
          <div class="games-page-toolbar">
            <div class="games-page-toolbar-left">
              <input
                class="games-page-search-input"
                type="text"
                placeholder="Search games"
              />
            </div>
            <div class="games-page-toolbar-right">
              <button class="games-page-icon-button" type="button">
                ⟳
              </button>
              <button class="games-page-icon-button" type="button">
                ⚙️
              </button>
            </div>
          </div>

          <div class="games-page-list-wrapper">
            <GameList />
          </div>
        </div>
      </section>
    </div>
  );
};

