import "./GamesPage.css";
import { GameList } from "../components/GameList";
import { GameFormModal } from "../components/GameFormModal";
import { gameStore } from "../store/game.store.ts";
import { createEffect, createSignal, Show } from "solid-js";
import type { CreateGameDto, UpdateGameDto } from "../types/game.types";

export const GamesPage = () => {
  const { state, actions } = gameStore;
  const [isCreateModalOpen, setIsCreateModalOpen] = createSignal(false);
  const [isEditModalOpen, setIsEditModalOpen] = createSignal(false);

  createEffect(() => {
    actions.loadGames();
  });

  const handleCreate = async (dto: CreateGameDto) => {
    try {
      await actions.createGame(dto);
      setIsCreateModalOpen(false);
    } catch (error) {
      // Ошибка будет обработана в модальном окне
      throw error;
    }
  };

  const handleUpdate = async (dto: UpdateGameDto) => {
    if (state.selectedGame) {
      try {
        await actions.updateGame(state.selectedGame.id, dto);
        setIsEditModalOpen(false);
        actions.setSelectedGame(null);
      } catch (error) {
        // Ошибка будет обработана в модальном окне
        throw error;
      }
    }
  };

  const handleDelete = async () => {
    if (state.selectedGame) {
      if (confirm(`Вы уверены, что хотите удалить игру "${state.selectedGame.title}"?`)) {
        await actions.deleteGame(state.selectedGame.id);
        actions.setSelectedGame(null);
      }
    }
  };

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
            <button
              textContent={"Обновить"}
              class="games-page-icon-button"
              type="button"
              onClick={() => {
                actions.loadGames();
              }}
            />
          </div>
        </div>

        <Show when={state.selectedGame}>
          <div class="games-page-actions-bar">
            <button
              class="games-page-action-button games-page-action-button--edit"
              type="button"
              onClick={() => setIsEditModalOpen(true)}
              disabled={state.isLoading}
            >
              Изменить
            </button>
            <button
              class="games-page-action-button games-page-action-button--delete"
              type="button"
              onClick={handleDelete}
              disabled={state.isLoading}
            >
              Удалить
            </button>
          </div>
        </Show>

        <div class="games-page-list-wrapper">
          <GameList state={state} actions={actions} />
        </div>

        <div class="games-page-create-button-container">
          <button
            class="games-page-create-button"
            type="button"
            onClick={() => setIsCreateModalOpen(true)}
            disabled={state.isLoading}
          >
            + Создать новую игру
          </button>
        </div>
      </div>

      <GameFormModal
        isOpen={isCreateModalOpen()}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={handleCreate}
        isLoading={state.isLoading}
      />

      <GameFormModal
        isOpen={isEditModalOpen()}
        game={state.selectedGame}
        onClose={() => setIsEditModalOpen(false)}
        onSubmit={handleUpdate}
        isLoading={state.isLoading}
      />
    </div>
  );
};

