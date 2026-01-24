import "./GamesPage.css";
import { GameList } from "../components/GameList";
import { GameFormModal } from "../components/GameFormModal";
import { Modal } from "../../../shared/components/modal/Modal.tsx";
import { gameStore } from "../store/game.store.ts";
import { onMount, createSignal, Show, createEffect } from "solid-js";
import type { CreateGameDto, UpdateGameDto } from "../types/game.types";
import { authStore } from "../../auth/store/auth.store";
import { unityApi, type UnityStatus } from "../api/unity.api.ts";

const UNITY_PATH_STORAGE_KEY = "unityExecutablePath";
const DEFAULT_UNITY_PATH = "/traffic_light_alarm_system/Система сигнализации светофоров.x86_64";

export const GamesPage = () => {
  const { state, actions } = gameStore;
  const [isCreateModalOpen, setIsCreateModalOpen] = createSignal(false);
  const [isEditModalOpen, setIsEditModalOpen] = createSignal(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = createSignal(false);
  const isAdmin = () => authStore.actions.isAdmin();
  const [isUnityRunning, setIsUnityRunning] = createSignal(false);
  const [unityPid, setUnityPid] = createSignal<number | null>(null);
  const [isUnityBusy, setIsUnityBusy] = createSignal(false);
  const [unityError, setUnityError] = createSignal<string | null>(null);
  const initialUnityExecutablePath = (() => {
    try {
      return localStorage.getItem(UNITY_PATH_STORAGE_KEY) || DEFAULT_UNITY_PATH;
    } catch {
      return DEFAULT_UNITY_PATH;
    }
  })();
  const [unityExecutablePath, setUnityExecutablePath] =
    createSignal<string>(initialUnityExecutablePath);

  createEffect(() => {
    // #region agent log
    fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:createEffect',message:'Unity UI state changed',data:{busy:isUnityBusy(),running:isUnityRunning(),pid:unityPid(),hasError:!!unityError()},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'G'})}).catch(()=>{});
    // #endregion
  });

  onMount(() => {
    actions.loadGames();
    unityApi
      .status()
      .then((s) => {
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:onMount',message:'Unity status on mount',data:{running:s.running,pid:s.pid},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'E'})}).catch(()=>{});
        // #endregion
        setIsUnityRunning(!!s.running);
        setUnityPid(s.pid);
      })
      .catch((e) => {
        const msg = e instanceof Error ? e.message : String(e);
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:onMount',message:'Unity status on mount error',data:{error:msg},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'F'})}).catch(()=>{});
        // #endregion
        // ignore: Unity feature not critical for list view
      });
  });

  const startUnity = async () => {
    setUnityError(null);
    setIsUnityBusy(true);
    try {
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:startUnity',message:'Starting Unity',data:{selectedGameId:state.selectedGame?.id ?? null,executablePath:unityExecutablePath()},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'A'})}).catch(()=>{});
      // #endregion
      const status: UnityStatus = await unityApi.start(unityExecutablePath());
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:startUnity',message:'Unity start result',data:{running:status.running,pid:status.pid},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'A'})}).catch(()=>{});
      // #endregion
      setIsUnityRunning(!!status.running);
      setUnityPid(status.pid);
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:startUnity',message:'Unity start error',data:{error:msg},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'B'})}).catch(()=>{});
      // #endregion
      setUnityError(msg);
      setIsUnityRunning(false);
      setUnityPid(null);
    } finally {
      setIsUnityBusy(false);
    }
  };

  const stopUnity = async () => {
    setUnityError(null);
    setIsUnityBusy(true);
    try {
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:stopUnity',message:'Stopping Unity',data:{pid:unityPid()},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'C'})}).catch(()=>{});
      // #endregion
      const status: UnityStatus = await unityApi.stop();
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:stopUnity',message:'Unity stop result',data:{running:status.running,pid:status.pid},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'C'})}).catch(()=>{});
      // #endregion
      setIsUnityRunning(false);
      setUnityPid(null);
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'GamesPage.tsx:stopUnity',message:'Unity stop error',data:{error:msg},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'D'})}).catch(()=>{});
      // #endregion
      setUnityError(msg);
    } finally {
      setIsUnityBusy(false);
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
                  disabled={state.isLoading}
                  onClick={() => {
                    startUnity();
                  }}
                >
                  Начать игру
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
        isOpen={isUnityRunning() || isUnityBusy() || !!unityError()}
        title={isUnityRunning() || isUnityBusy() ? "Игра запущена" : "Не удалось запустить игру"}
        onClose={() => {
          if (!isUnityRunning() && !isUnityBusy()) setUnityError(null);
        }}
        showCloseButton={!isUnityRunning() && !isUnityBusy()}
        closeOnOverlayClick={!isUnityRunning() && !isUnityBusy()}
        footer={
          <>
            <Show
              when={isUnityRunning() || isUnityBusy()}
              fallback={
                <>
                  <button
                    class="modal-btn"
                    onClick={() => setUnityError(null)}
                    disabled={isUnityBusy()}
                  >
                    Закрыть
                  </button>
                  <button class="modal-btn primary" onClick={startUnity} disabled={isUnityBusy()}>
                    {isUnityBusy() ? "Запуск..." : "Попробовать снова"}
                  </button>
                </>
              }
            >
              <button class="modal-btn danger" onClick={stopUnity} disabled={isUnityBusy()}>
                {isUnityBusy() ? "Остановка..." : "выйти"}
              </button>
            </Show>
          </>
        }
      >
        <Show
          when={isUnityRunning() || isUnityBusy()}
          fallback={
            <>
              <p style={{ "margin-bottom": "0.5rem" }}>{unityError()}</p>
              <p style={{ "margin-bottom": "0.35rem", color: "#6b7280", "font-size": "0.9rem" }}>
                Путь к исполняемому файлу Unity:
              </p>
              <input
                class="games-page-search-input"
                type="text"
                value={unityExecutablePath()}
                onInput={(e) => {
                  const v = e.target.value;
                  setUnityExecutablePath(v);
                  try {
                    localStorage.setItem(UNITY_PATH_STORAGE_KEY, v);
                  } catch {}
                }}
              />
            </>
          }
        >
          <p style={{ "margin-bottom": "0.5rem" }}>
            Пока игра запущена, основное приложение недоступно.
          </p>
          <Show when={unityPid()}>
            <p style={{ color: "#6b7280", "font-size": "0.9rem" }}>PID: {unityPid()}</p>
          </Show>
        </Show>
      </Modal>
    </div>
  );
};

