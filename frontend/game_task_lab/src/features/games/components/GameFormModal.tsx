import { createSignal, Show, createEffect } from "solid-js";
import type { CreateGameDto, GameDto, UpdateGameDto } from "../types/game.types";
import type { GameStore } from "../store/game.store";
import { Modal } from "../../../shared/components/modal/Modal.tsx";
import { FlatpickrInput } from "../../../shared/components/flatpickr/FlatpickrInput.tsx";

interface GameFormModalProps {
  isOpen: boolean;
  game?: GameDto | null;
  onClose: () => void;
  onSubmit: (dto: CreateGameDto | UpdateGameDto) => Promise<void>;
  isLoading?: boolean;
  gameStore?: GameStore;
}

export const GameFormModal = (props: GameFormModalProps) => {
  const [title, setTitle] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [releaseDate, setReleaseDate] = createSignal("");
  const [error, setError] = createSignal<string | null>(null);

  createEffect(() => {
    if (props.isOpen) {
      if (props.gameStore) {
        const consumedError = props.gameStore.actions.consumeError();
        if (consumedError) {
          setError(consumedError);
        } else {
          setError(null);
        }
      } else {
        setError(null);
      }

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
      if (props.gameStore) {
        const consumedError = props.gameStore.actions.consumeError();
        setError(consumedError || (err instanceof Error ? err.message : "Произошла ошибка при сохранении"));
      } else {
        setError(err instanceof Error ? err.message : "Произошла ошибка при сохранении");
      }
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
                <Show when={props.isOpen}>
                  <FlatpickrInput
                    value={releaseDate()}
                    onChange={(date) => setReleaseDate(date)}
                    required
                    disabled={props.isLoading}
                    placeholder="Выберите дату"
                    style={{
                      width: "100%",
                      padding: "0.5rem",
                      "border-radius": "6px",
                      border: "1px solid #d1d5db",
                      "background-color": props.isLoading ? "#f3f4f6" : "#ffffff",
                      color: "#111827",
                      cursor: props.isLoading ? "not-allowed" : "pointer",
                    }}
                  />
                </Show>
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
