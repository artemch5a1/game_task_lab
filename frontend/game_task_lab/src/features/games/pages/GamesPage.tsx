import "./GamesPage.css";
import { GameList } from "../components/GameList";
import { GameFormModal } from "../components/GameFormModal";
import { Modal } from "../../../shared/components/modal/Modal.tsx";
import { gameStore } from "../store/game.store.ts";
import { onMount, createSignal, Show } from "solid-js";
import type { CreateGameDto, UpdateGameDto } from "../types/game.types";
import { authStore } from "../../auth/store/auth.store";
import { webglApi, type WebglStatus } from "../api/webgl.api.ts";

export const GamesPage = () => {
  const { state, actions } = gameStore;
  const [isCreateModalOpen, setIsCreateModalOpen] = createSignal(false);
  const [isEditModalOpen, setIsEditModalOpen] = createSignal(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = createSignal(false);
  const isAdmin = () => authStore.actions.isAdmin();
  const [isGameBusy, setIsGameBusy] = createSignal(false);
  const [gameError, setGameError] = createSignal<string | null>(null);

  onMount(() => {
    actions.loadGames();
  });

  const startGame = async () => {
    setGameError(null);
    setIsGameBusy(true);
    try {
      const status: WebglStatus = await webglApi.start();
      if (!status.url) throw new Error("WebGL server did not return a URL");
      // Load the WebGL build in the same Tauri window (no separate OS window).
      const returnUrl = window.location.href;
      const url = `${status.url}?return=${encodeURIComponent(returnUrl)}`;
      window.location.replace(url);
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      setGameError(msg);
    } finally {
      setIsGameBusy(false);
    }
  };

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
        throw error;
      }
    }
  };

  const handleDeleteClick = () => {
    if (state.selectedGame) {
      setIsDeleteModalOpen(true);
    }
  };

  const handleDeleteConfirm = async () => {
    if (state.selectedGame) {
      try {
        await actions.deleteGame(state.selectedGame.id);
        actions.setSelectedGame(null);
        setIsDeleteModalOpen(false);
      } catch (error) {
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
            <Show
              when={isAdmin()}
              fallback={
                <button
                  class="games-page-action-button games-page-action-button--start"
                  type="button"
                  disabled={state.isLoading || isGameBusy()}
                  onClick={() => {
                    startGame();
                  }}
                >
                  {isGameBusy() ? "Запуск..." : "Начать игру"}
                </button>
              }
            >
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
                onClick={handleDeleteClick}
                disabled={state.isLoading}
              >
                Удалить
              </button>
            </Show>
          </div>
        </Show>

        <div class="games-page-list-wrapper">
          <GameList state={state} actions={actions} />
        </div>

        <Show when={isAdmin()}>
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
        </Show>
      </div>

      <GameFormModal
        isOpen={isCreateModalOpen()}
        onClose={() => {
          actions.clearError();
          setIsCreateModalOpen(false);
        }}
        onSubmit={handleCreate}
        isLoading={state.isLoading}
        gameStore={gameStore}
      />

      <GameFormModal
        isOpen={isEditModalOpen()}
        game={state.selectedGame}
        onClose={() => {
          actions.clearError();
          setIsEditModalOpen(false);
        }}
        onSubmit={handleUpdate}
        isLoading={state.isLoading}
        gameStore={gameStore}
      />

      <Modal
        isOpen={isDeleteModalOpen()}
        title="Подтверждение удаления"
        onClose={() => setIsDeleteModalOpen(false)}
        footer={
          <>
            <button
              class="modal-btn"
              onClick={() => setIsDeleteModalOpen(false)}
              disabled={state.isLoading}
            >
              Отмена
            </button>
            <button
              class="modal-btn danger"
              onClick={handleDeleteConfirm}
              disabled={state.isLoading}
            >
              {state.isLoading ? "Удаление..." : "Удалить"}
            </button>
          </>
        }
      >
        <p>
          Вы уверены, что хотите удалить игру{" "}
          <strong>{state.selectedGame?.title}</strong>?
        </p>
        <p style={{ "margin-top": "0.5rem", color: "#6b7280", "font-size": "0.9rem" }}>
          Это действие нельзя отменить.
        </p>
      </Modal>

      <Modal
        isOpen={isGameBusy() || !!gameError()}
        title={isGameBusy() ? "Запуск игры" : "Не удалось запустить игру"}
        onClose={() => {
          if (!isGameBusy()) setGameError(null);
        }}
        showCloseButton={!isGameBusy()}
        closeOnOverlayClick={!isGameBusy()}
        footer={
          <>
            <button class="modal-btn" onClick={() => setGameError(null)} disabled={isGameBusy()}>
              Закрыть
            </button>
            <button class="modal-btn primary" onClick={startGame} disabled={isGameBusy()}>
              {isGameBusy() ? "Запуск..." : "Попробовать снова"}
            </button>
          </>
        }
      >
        <Show when={isGameBusy()} fallback={<p style={{ "margin-bottom": "0.5rem" }}>{gameError()}</p>}>
          <p style={{ "margin-bottom": "0.5rem" }}>
            Игра запускается. Окно приложения будет переключено на WebGL-сборку.
          </p>
        </Show>
      </Modal>
    </div>
  );
};

