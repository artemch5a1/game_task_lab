import "./GamesPage.css";
import { GameList } from "../components/GameList";

export const GamesPage = () => {
  return (
    <div class="games-page">
      <h1 class="games-page-title">Game List</h1>

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
    </div>
  );
};

