// features/games/components/GameList.tsx
import { createEffect } from "solid-js";
import {GameStore} from "../store/game.store";
import { Table, type TableColumn } from "../../../shared/components/table/Table.tsx";
import { Modal } from "../../../shared/components/modal/Modal.tsx";
import "./GameList.css";

const formatDate = (dateString: string): string => {
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString("ru-RU", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  } catch {
    return dateString;
  }
};



export const GameList = (gameStore: GameStore) => {
  const { state, actions } = gameStore;

  createEffect(() => {
    actions.loadGames();
  });

  const columns: TableColumn[] = [
    {
      key: "title",
      header: "Название",
    },
    {
      key: "description",
      header: "Описание",
    },
    {
      key: "releaseDate",
      header: "Дата релиза",
      render: (game) => formatDate(game.releaseDate),
    },
  ];
  const closeErrorModal = () => {
    gameStore.actions.setErrorNull();
  };

  return (
    <div class="game-list-container">
      <Modal
        isOpen={!!state.error}
        title="Ошибка"
        onClose={closeErrorModal}
        footer={
          <>
            <button
              class="modal-btn primary"
              onClick={() => {
                closeErrorModal();
                actions.loadGames();
              }}
            >
              Попробовать снова
            </button>
            <button
              class="modal-btn"
              onClick={closeErrorModal}
            >
              Закрыть
            </button>
          </>
        }
      >
        <p>{state.error}</p>
      </Modal>

      <Table
        columns={columns}
        data={state.games}
        isLoading={state.isLoading}
        emptyText="Игры не найдены"
        selectedRowId={state.selectedGame?.id ?? null}
        onRowClick={(game) => actions.setSelectedGame(game)}
      />
    </div>
  );
};