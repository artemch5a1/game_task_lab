import { createSignal, Show, createEffect } from "solid-js";
import type { CreateGameDto, GameDto, UpdateGameDto } from "../types/game.types";
import type { GameStore } from "../store/game.store";
import { Modal } from "../../../shared/components/modal/Modal.tsx";
import { FlatpickrInput } from "../../../shared/components/flatpickr/FlatpickrInput.tsx";
import { genreApi } from "../api/genre.api";
import type { GenreDto } from "../types/genre.types";

interface GameFormModalProps {
  isOpen: boolean;
  game?: GameDto | null;
  onClose: () => void;
  onSubmit: (dto: any) => Promise<void>;
  isLoading?: boolean;
  gameStore?: GameStore;
}

export const GameFormModal = (props: GameFormModalProps) => {
  const [title, setTitle] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [releaseDate, setReleaseDate] = createSignal("");
  const [genreId, setGenreId] = createSignal<string>("");
  const [genres, setGenres] = createSignal<GenreDto[]>([]);
  const [genresLoading, setGenresLoading] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  let genreSelectRef: HTMLSelectElement | undefined;
  const [genreSelectEl, setGenreSelectEl] = createSignal<HTMLSelectElement | null>(null);
  const [genreInitializedKey, setGenreInitializedKey] = createSignal<string | null>(null);
  const [genreDirty, setGenreDirty] = createSignal(false);
  let genreSelectMountSeq = 0;

  createEffect(() => {
    if (!props.isOpen) return;
    if (!genreSelectRef) return;
    if (genres().length === 0) return;

    // #region agent log
    fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'E',location:'GameFormModal.tsx:domSelect:stateVsDom',message:'compare genreId signal vs select DOM',data:{gameId:props.game?.id??null,signalGenreId:genreId(),domValue:genreSelectRef.value,selectedIndex:genreSelectRef.selectedIndex,selectedOptionValue:genreSelectRef.selectedOptions?.[0]?.value??null,optionsCount:genreSelectRef.options?.length??null},timestamp:Date.now()})}).catch(()=>{});
    // #endregion
  });

  // На некоторых ре-монтажах модалки появляется новый <select>, который сбрасывается на placeholder.
  // Этот эффект держит DOM в соответствии с состоянием.
  createEffect(() => {
    if (!props.isOpen) return;
    if (genres().length === 0) return;
    const el = genreSelectEl();
    if (!el) return;

    const desired = props.game
      ? (genreDirty() ? genreId() : props.game.genreId)
      : genreId();

    // #region agent log
    fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'post-fix2',hypothesisId:'G',location:'GameFormModal.tsx:domSelect:enforce',message:'enforce select DOM value',data:{gameId:props.game?.id??null,genreDirty:genreDirty(),signalGenreId:genreId(),desired,domValueBefore:el.value,optionsCount:el.options?.length??null},timestamp:Date.now()})}).catch(()=>{});
    // #endregion

    if (desired && el.value !== desired) {
      el.value = desired;
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'post-fix2',hypothesisId:'G',location:'GameFormModal.tsx:domSelect:enforceApplied',message:'applied enforce select DOM value',data:{gameId:props.game?.id??null,desired,domValueAfter:el.value,selectedIndex:el.selectedIndex,selectedOptionValue:el.selectedOptions?.[0]?.value??null},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
    }
  });

  // Когда жанры догружаются, гарантируем, что select отобразит текущий жанр при редактировании
  // (если value был проставлен до появления options, браузер может не выбрать нужный option автоматически).
  createEffect(() => {
    // #region agent log
    fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'A',location:'GameFormModal.tsx:syncEffect:entry',message:'syncEffect entry',data:{isOpen:props.isOpen,gameId:props.game?.id??null,gameGenreId:props.game?.genreId??null,genreId:genreId(),genresLen:genres().length},timestamp:Date.now()})}).catch(()=>{});
    // #endregion
    if (!props.isOpen) return;

    const list = genres();
    if (list.length === 0) return;

    // Делаем эффект зависимым от реального DOM-элемента select (реф сам по себе не реактивный).
    const el = genreSelectEl();

    const key = props.game ? `edit:${props.game.id}` : "create";
    if (genreInitializedKey() === key) return;

    if (props.game) {
        console.error('ыыыыыыыыыыыыыыыыы1');

      const current = props.game.genreId;
      // Если пользователь уже менял жанр — не перетираем его выбором из props.game.
      if (genreDirty()) return;

      if (list.some((g) => g.id === current)) {
          console.error('ыыыыыыыыыыыыыыыыы2 [current]', current);
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'F',location:'GameFormModal.tsx:syncEffect:willForceToGameGenre',message:'syncEffect in edit will force genreId to props.game.genreId',data:{gameId:props.game.id,propsGameGenreId:props.game.genreId,currentSignalGenreId:genreId(),domValue:genreSelectRef?.value??null},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'C',location:'GameFormModal.tsx:syncEffect:setGenreId',message:'syncEffect setGenreId from props.game.genreId',data:{gameId:props.game.id,from:genreId(),to:current,genresLen:list.length},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
        setGenreId(current);
        // Принудительно синхронизируем DOM select (на некоторых ре-монтажах модалки браузер оставляет placeholder)
        if (el) {
          el.value = current;
          // #region agent log
          fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'E',location:'GameFormModal.tsx:syncEffect:forceDomValue',message:'forced select DOM value after setGenreId',data:{gameId:props.game.id,forcedValue:current,domValue:el.value,selectedIndex:el.selectedIndex,selectedOptionValue:el.selectedOptions?.[0]?.value??null},timestamp:Date.now()})}).catch(()=>{});
          // #endregion
        }
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'A',location:'GameFormModal.tsx:syncEffect:afterSet',message:'syncEffect after setGenreId',data:{genreId:genreId()},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
      }

      // Инициализацию считаем завершённой только когда select реально смонтирован.
      if (!el) {
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'A',location:'GameFormModal.tsx:syncEffect:skipInitNoEl',message:'skip initKey set because select element not mounted',data:{key,gameId:props.game.id},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
        return;
      }

      setGenreInitializedKey(key);
      return;
    }

    // Для создания можно проставить первый жанр по умолчанию, чтобы форма была валидной
    if (!genreId()) {
        console.error('ыыыыыыыыыыыыыыыыы3');
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'D',location:'GameFormModal.tsx:syncEffect:defaultCreateGenre',message:'default genreId for create',data:{to:list[0]?.id??null,genresLen:list.length},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
      setGenreId(list[0]!.id);
      if (el) {
        el.value = list[0]!.id;
        // #region agent log
        fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'E',location:'GameFormModal.tsx:syncEffect:forceDomValueCreate',message:'forced select DOM value for create default',data:{forcedValue:list[0]!.id,domValue:el.value,selectedIndex:el.selectedIndex,selectedOptionValue:el.selectedOptions?.[0]?.value??null},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
      }
    }
    if (!el) return;
    setGenreInitializedKey(key);
  });

  createEffect(() => {
    if (props.isOpen) {
      setGenreInitializedKey(null);
      setGenreDirty(false);
      // #region agent log
      fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'B',location:'GameFormModal.tsx:openEffect:entry',message:'openEffect entry',data:{isOpen:props.isOpen,gameId:props.game?.id??null,gameGenreId:props.game?.genreId??null},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
      // Каждый раз при открытии актуализируем жанры (backend мог измениться)
      setGenresLoading(true);
      genreApi
        .getAllGenres()
        .then((list) => {
          setGenres(list);
          // #region agent log
          fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'B',location:'GameFormModal.tsx:openEffect:genresLoaded',message:'genres loaded',data:{genresLen:list.length,firstId:list[0]?.id??null,firstTitle:list[0]?.title??null},timestamp:Date.now()})}).catch(()=>{});
          // #endregion
        })
        .catch((err) => {
          setGenres([]);
          setError(err instanceof Error ? err.message : "Не удалось загрузить жанры");
          // #region agent log
          fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'B',location:'GameFormModal.tsx:openEffect:genresError',message:'genres load error',data:{error:err instanceof Error ? err.message : String(err)},timestamp:Date.now()})}).catch(()=>{});
          // #endregion
        })
        .finally(() => {
          setGenresLoading(false);
        });

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
        setGenreId(props.game.genreId);
        setReleaseDate(
          props.game.releaseDate
            ? new Date(props.game.releaseDate).toISOString().split("T")[0]
            : ""
        );
      } else {
        setTitle("");
        setDescription("");
        setReleaseDate("");
        setGenreId("");
      }
    }
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError(null);
    
    try {
      const dto: CreateGameDto | UpdateGameDto = {
        title: title(),
        description: description().trim() || undefined,
        releaseDate: new Date(releaseDate()).toISOString(),
        genreId: genreId(),
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

              <div style={{ "margin-bottom": "1rem" }}>
                <label style={{ display: "block", "margin-bottom": "0.5rem", "font-weight": "500" }}>
                  Жанр *
                </label>
                <select
                  ref={(el) => {
                    genreSelectMountSeq += 1;
                    genreSelectRef = el;
                    setGenreSelectEl(el);
                    // #region agent log
                    fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'post-fix2',hypothesisId:'G',location:'GameFormModal.tsx:select:ref',message:'select ref assigned',data:{mountSeq:genreSelectMountSeq,gameId:props.game?.id??null,domValue:el.value,optionsCount:el.options?.length??null},timestamp:Date.now()})}).catch(()=>{});
                    // #endregion
                  }}
                  value={genreId()}
                  onChange={(e) => {
                    const next = e.currentTarget.value;
                    // #region agent log
                    fetch('http://127.0.0.1:7243/ingest/dbfb5b5d-abe1-4252-bd32-bcd21c6938c0',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'F',location:'GameFormModal.tsx:select:onChange',message:'user changed genre select',data:{gameId:props.game?.id??null,fromSignal:genreId(),to:next,domValue:e.currentTarget.value,selectedIndex:e.currentTarget.selectedIndex,selectedOptionText:e.currentTarget.selectedOptions?.[0]?.text??null},timestamp:Date.now()})}).catch(()=>{});
                    // #endregion
                    setGenreDirty(true);
                    setGenreId(next);
                  }}
                  required
                  disabled={props.isLoading || genresLoading()}
                  style={{
                    width: "100%",
                    padding: "0.5rem",
                    "border-radius": "6px",
                    border: "1px solid #d1d5db",
                    "background-color":
                      props.isLoading || genresLoading() ? "#f3f4f6" : "#ffffff",
                    color: "#111827",
                    cursor: props.isLoading || genresLoading() ? "not-allowed" : "pointer",
                  }}
                >
                  <option value="" disabled>
                    {genresLoading() ? "Загрузка жанров..." : "Выберите жанр"}
                  </option>
                  {genres().map((g) => (
                    <option value={g.id}>{g.title}</option>
                  ))}
                </select>
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
