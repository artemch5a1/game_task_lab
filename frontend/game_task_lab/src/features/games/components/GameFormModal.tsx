import { createSignal, Show, createEffect } from "solid-js";
import type { CreateGameDto, GameDto, UpdateGameDto } from "../types/game.types";
import { Modal } from "../../../shared/components/modal/Modal.tsx";

interface GameFormModalProps {
  isOpen: boolean;
  game?: GameDto | null;
  onClose: () => void;
  onSubmit: (dto: CreateGameDto | UpdateGameDto) => Promise<void>;
  isLoading?: boolean;
}

export const GameFormModal = (props: GameFormModalProps) => {
  const [title, setTitle] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [releaseDate, setReleaseDate] = createSignal("");
  const [error, setError] = createSignal<string | null>(null);

  createEffect(() => {
    if (props.isOpen) {
      setError(null);
      if (props.game) {
        setTitle(props.game.title);
        setDescription(props.game.description ?? "");
        setReleaseDate(
          props.game.releaseDate
            ? new Date(props.game.releaseDate).toISOString().split("T")[0]
            : ""
        );
      } else {
        setTitle("");
        setDescription("");
        setReleaseDate("");
      }
    }
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError(null);
    
    try {
      // Автоматически генерируем UUID для genreId только при создании новой игры
      // При обновлении используем существующий genreId
      const genreId = props.game ? props.game.genreId : crypto.randomUUID();
      
      const dto: CreateGameDto | UpdateGameDto = {
        title: title(),
        description: description().trim() || undefined,
        releaseDate: new Date(releaseDate()).toISOString(),
        genreId: genreId,
        ...(props.game ? { id: props.game.id } : {}),
      };
      await props.onSubmit(dto);
      props.onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Произошла ошибка при сохранении");
    }
  };

  const handleClose = () => {
    setError(null);
    props.onClose();
  };

  const formId = "game-form-modal";

  return (
    <Modal
      isOpen={props.isOpen}
      title={props.game ? "Редактировать игру" : "Создать игру"}
      onClose={handleClose}
      footer={
        <>
          <button
            type="button"
            class="modal-btn"
            onClick={handleClose}
            disabled={props.isLoading}
          >
            Отмена
          </button>
          <button
            type="submit"
            form={formId}
            class="modal-btn primary"
            disabled={props.isLoading}
          >
            {props.isLoading ? "Сохранение..." : props.game ? "Сохранить" : "Создать"}
          </button>
        </>
      }
    >
      <form id={formId} onSubmit={handleSubmit}>
              <div style={{ "margin-bottom": "1rem" }}>
                <label style={{ display: "block", "margin-bottom": "0.5rem", "font-weight": "500" }}>
                  Название *
                </label>
                <input
                  type="text"
                  value={title()}
                  onInput={(e) => setTitle(e.currentTarget.value)}
                  required
                  disabled={props.isLoading}
                  style={{
                    width: "100%",
                    padding: "0.5rem",
                    "border-radius": "6px",
                    border: "1px solid #d1d5db",
                  }}
                />
              </div>

              <div style={{ "margin-bottom": "1rem" }}>
                <label style={{ display: "block", "margin-bottom": "0.5rem", "font-weight": "500" }}>
                  Описание
                </label>
                <textarea
                  value={description()}
                  onInput={(e) => setDescription(e.currentTarget.value)}
                  disabled={props.isLoading}
                  rows={4}
                  style={{
                    width: "100%",
                    padding: "0.5rem",
                    "border-radius": "6px",
                    border: "1px solid #d1d5db",
                    "font-family": "inherit",
                  }}
                />
              </div>

              <div style={{ "margin-bottom": "1rem" }}>
                <label style={{ display: "block", "margin-bottom": "0.5rem", "font-weight": "500" }}>
                  Дата релиза *
                </label>
                <input
                  type="date"
                  value={releaseDate()}
                  onInput={(e) => setReleaseDate(e.currentTarget.value)}
                  onBlur={(e) => {
                    // Позволяем закрыть календарь при потере фокуса
                    e.currentTarget.blur();
                  }}
                  onKeyDown={(e) => {
                    // Позволяем закрыть календарь по Escape
                    if (e.key === "Escape") {
                      e.currentTarget.blur();
                    }
                  }}
                  required
                  disabled={props.isLoading}
                  style={{
                    width: "100%",
                    padding: "0.5rem",
                    "border-radius": "6px",
                    border: "1px solid #d1d5db",
                  }}
                />
              </div>

              <Show when={error()}>
                <div style={{
                  padding: "0 0 10px",
                  color: "#dc3545",
                  "font-size": "0.9rem",
                  "text-align": "center",
                }}>
                  {error()}
                </div>
              </Show>
            </form>
    </Modal>
  );
};
