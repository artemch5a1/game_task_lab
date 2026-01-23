// features/games/components/GameList.tsx
import { createEffect, Show } from "solid-js";
import { gameStore } from "../store/game.store";
import { Table, type TableColumn } from "../../../components/Table";
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
  {
    key: "actions",
    header: "Действия",
    align: "right",
    render: (game) => (
      <button
        class="delete-button"
        onClick={() => {
          if (
            confirm(`Вы уверены, что хотите удалить игру "${game.title}"?`)
          ) {
            gameStore.actions.deleteGame(game.id);
          }
        }}
        disabled={gameStore.state.isLoading}
      >
        Удалить
      </button>
    ),
  },
];

export const GameList = () => {
  const { state, actions } = gameStore;

  createEffect(() => {
    actions.loadGames();
  });

  return (
    <div class="game-list-container">
      <Show when={state.error}>
        <div class="error-container">
          <p class="error-message">Ошибка: {state.error}</p>
          <button class="retry-button" onClick={() => actions.loadGames()}>
            Попробовать снова
          </button>
        </div>
      </Show>

      <Table
        columns={columns}
        data={state.games}
        isLoading={state.isLoading}
        emptyText="Игры не найдены"
      />
    </div>
  );
};