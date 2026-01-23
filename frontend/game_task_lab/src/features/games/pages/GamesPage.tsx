import "./GamesPage.css";
import { GameList } from "../components/GameList";
import {gameStore} from "../store/game.store.ts";
import {createEffect} from "solid-js";

export const GamesPage = () => {

  const { state, actions } = gameStore;

  createEffect(() => {
    actions.loadGames();
  });

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
              onInput={(e) => {
                const searchString = e.target.value;
                actions.filterGames(searchString);
              }}
            />
          </div>
          <div class="games-page-toolbar-right">
            <button textContent={"Обновить"} class="games-page-icon-button" type="button" onClick={() => { actions.loadGames(); }}/>
          </div>
        </div>

        <div class="games-page-list-wrapper">
          <GameList state={state} actions={actions} />
        </div>
      </div>
    </div>
  );
};

